import pandas as pd
from sklearn import model_selection, metrics, preprocessing
import torch
import torch.nn as nn
import torch.optim as optim
import sqlite3
from torch.utils.data import DataLoader, Dataset

device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

db = sqlite3.connect('db/westflix.db')
cur = db.cursor()
res = cur.execute("SELECT * FROM RecommenderFunction")
ratings = res.fetchall()

cur = db.cursor()
res = cur.execute("SELECT * FROM User")
users = res.fetchall()

df = pd.DataFrame.from_records(ratings, columns=['movieId', 'userId', 'ratings'])

# print(df.info)
# print(df.movieId.nunique())
# print(df.userId.nunique())
# print(df.ratings.value_counts())
# print(df.shape)


class ClemflixDataset:
    def __init__(self, users, movies, ratings):
        self.users = users
        self.movies = movies
        self.ratings = ratings
    # endDef

    def __len__(self):
        return len(self.users)
    # endDef

    def __getitem__(self, item):
        users = self.users[item]
        movies = self.movies[item]
        ratings = self.ratings[item]

        return {
            "users": torch.tensor(users, dtype=torch.int64),
            "movies": torch.tensor(movies, dtype=torch.int64),
            "ratings": torch.tensor(ratings, dtype=torch.float64),
        }
    # endDef
# endClass


class RecSysModel(nn.Module):
    def __init__(self, n_users, n_movies):
        super().__init__()

        self.user_embed = nn.Embedding(n_users, 32)
        self.movie_embed = nn.Embedding(n_movies, 32)

        self.out = nn.Linear(64, 1)
    # endDef

    def forward(self, users, movies, ratings=None):
        user_embeds = self.user_embed(users)
        movie_embeds = self.movie_embed(movies)
        output = torch.cat([user_embeds, movie_embeds], dim=1)

        output = self.out(output)
        return output
    # endDef
# endClass


# Preprocess movie and user id's so that there are index of of bound errors
lbl_user = preprocessing.LabelEncoder()
lbl_movie = preprocessing.LabelEncoder()
df.userId = lbl_user.fit_transform(df.userId.values)
df.movieId = lbl_movie.fit_transform(df.movieId.values)

df_train, df_valid = model_selection.train_test_split(
    df, test_size=(0.1), random_state=42
)

train_dataset = ClemflixDataset(
    users=df_train.userId.values,
    movies=df_train.movieId.values,
    ratings=df_train.ratings.values
)

valid_dataset = ClemflixDataset(
    users=df_train.userId.values,
    movies=df_train.movieId.values,
    ratings=df_train.ratings.values
)

train_loader = DataLoader(dataset=train_dataset, batch_size=4, shuffle=True, num_workers=2)
validation_loader = DataLoader(dataset=valid_dataset, batch_size=4, shuffle=True, num_workers=2)

dataIter = iter(train_loader)
dataloader_data = next(dataIter)
print(dataloader_data)

model = RecSysModel(
    n_users=len(df.userId),
    n_movies=len(df.movieId),
).to(device)


optimizer = torch.optim.Adam(model.parameters())
sch = torch.optim.lr_scheduler.StepLR(optimizer, step_size=3, gamma=0.7)

loss_func = nn.MSELoss()

print(len(lbl_user.classes_))
print(len(lbl_movie.classes_))
print(df.movieId.max())
print(len(train_dataset))

print(dataloader_data['users'])
print(dataloader_data['users'].size())
print(dataloader_data['movies'])
print(dataloader_data['movies'].size())

user_embed = nn.Embedding(len(lbl_user.classes_), 32)
movie_embed = nn.Embedding(len(lbl_movie.classes_), 32)

out = nn.Linear(64, 1)

user_embeds = user_embed(dataloader_data['users'])
movie_embeds = movie_embed(dataloader_data['movies'])

print(f"user embeds {user_embeds.size()}")
print(f"user embeds {user_embeds}")
print(f"movie embeds {movie_embeds.size()}")
print(f"movie embeds {movie_embeds}")

output = torch.cat([user_embeds, movie_embeds], dim=1)
print(f"output {output.size()}")
print(f"output {output}")
output = out(output)
print(f"output {output}")

with torch.no_grad():
    model_output = model(dataloader_data['users'], dataloader_data['movies'])
    print(f"model Output: {model_output}, size: {model_output.size()}")
# endWith

# Reshape ratings
rating = dataloader_data['ratings']

epochs = 200
total_loss = 0
plot_steps, print_steps = 5000, 5000
step_cnt = 0
all_losses_list = []

model.train()
for epoch_i in range(epochs):
    for i, train_data in enumerate(train_loader):
        output = model(train_data['users'], train_data['movies'])
        # .view(4,-1) is to reshape the ratingto match the shape of the modek output which is 4x1
        rating = train_data['ratings'].view(train_data['ratings'].size(0), -1).to(torch.float32)
        loss = loss_func(output, rating)
        total_loss = total_loss + loss.sum().item()
        optimizer.zero_grad()
        loss.backward()
        optimizer.step()
        step_cnt = step_cnt + len(train_data['users'])
        if step_cnt % plot_steps == 0:
            avg_loss = total_loss/(len(train_data['users']) * plot_steps)
            print(f"epoch {epoch_i} loss at step: {step_cnt} is {avg_loss}")
            all_losses_list.append(avg_loss)
            total_loss = 0  # reset total loss
        # endIf
    # endFor
# endFor

from sklearn.metrics import mean_squared_error

model_output_list = []
target_rating_list = []

model.eval()

with torch.no_grad():
    for i, batched_data in enumerate(validation_loader):
        model_output = model(batched_data['users'], batched_data['movies'])
        model_output_list.append(model_output.sum().item() / len(batched_data['users']))
        target_rating = batched_data['ratings']
        target_rating_list.append(target_rating.sum().item() / len(batched_data['users']))
        print(f"model_output: {model_output}, target_rating: {target_rating}")
    # endFor
# endWith

rms = mean_squared_error(target_rating_list, model_output_list, squared=False)
print(f"rms: {rms}")

from collections import defaultdict
user_est_true = defaultdict(list)

with torch.no_grad():
    for i, batched_data in enumerate(validation_loader):
        users = batched_data['users']
        movies = batched_data['movies']
        ratings = batched_data['ratings']
        model_output = model(batched_data['users'], batched_data['movies'])
        for i in range(len(users)):
            user_id = users[i].item()
            movie_id = movies[i].item()
            pred_rating = model_output[i][0].item()
            true_rating = ratings[i].item()
            # print(f"{user_id}, {movie_id}, {pred_rating}, {true_rating}")
            user_est_true[user_id].append({pred_rating, true_rating})
        # endFor
    # endFor
# endWith

with torch.no_grad():
    precisions = dict()
    recalls = dict()
    k = 10
    threshold = 3.5

    for uid, user_ratings in user_est_true.items():
        user_ratings.sort(key=lambda x: x[0], reverse=True)
        n_rel = sum((true_r >= threshold) for (_, true_r) in user_ratings)
        n_rec_k = sum((est >= threshold) for (est, _) in user_ratings[k])
        n_rel_and_rec_k = sum(((true_r >= threshold) and (est >= threshold)) for (est, true_r) in user_ratings[k])
        print(f"uid: {uid}, n_rel {n_rel}, n_rec_k: {n_rec_k}, n_rel and n_rec_k, {n_rel_and_rec_k}")

        precisions[uid] = n_rel_and_rec_k / n_rec_k if n_rec_k != 0 else 0
        recalls[uid] = n_rel_and_rec_k / n_rel if n_rel != 0 else 0
    # endFor
# endWith
