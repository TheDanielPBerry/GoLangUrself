import pandas as pd
import numpy as np
from sklearn import model_selection, metrics, preprocessing
import torch
import torch.nn as nn
import torch.optim as optim
import sqlite3
from torch.utils.data import DataLoader, Dataset


class ClemflixDataset(Dataset):
    def __init__(self, ratings):
        self.ratings = ratings
    # endDef

    def __len__(self):
        return len(self.ratings)
    # endDef

    def __getitem__(self, idx):
        movie_id = self.ratings[idx][0]
        user_id = self.ratings[idx][1]
        rating = self.ratings[idx][2]
        # user_id = self.ratings.iloc[idx]['userId']
        # movie_id = self.ratings.iloc[idx]['movieId']
        # rating = self.ratings.iloc[idx]['rating']
        return user_id, movie_id, rating
    # endDef
# endClass


class MatrixFactorization(nn.Module):
    def __init__(self, num_users, num_movies, embedding_dim):
        super(MatrixFactorization, self).__init__()
        self.user_embeddings = nn.Embedding(num_users, embedding_dim)
        self.movie_embeddings = nn.Embedding(num_movies, embedding_dim)
    # endDef

    def forward(self, user_id, movie_id):
        user_emb = self.user_embeddings(user_id)
        movie_emb = self.movie_embeddings(movie_id)
        return torch.matmul(user_emb, torch.transpose(movie_emb, 0, 1))
    # endDef
# endClass


db = sqlite3.connect('db/westflix.db')
cur = db.cursor()
res = cur.execute("SELECT * FROM RecommenderFunction")
ratings = res.fetchall()

cur = db.cursor()
res = cur.execute("SELECT MAX(UserId) FROM User")
userCount = res.fetchall()

num_users = userCount[0][0]+1
num_movies = 101
embedding_size = 32

dataset = ClemflixDataset(ratings)
dataloader = DataLoader(dataset, batch_size=32)

model = MatrixFactorization(num_users, num_movies, embedding_size)
criterion = nn.MSELoss()
optimizer = optim.Adam(model.parameters())

for epoch in range(500):
    for user_id, movie_id, rating in dataloader:
        optimizer.zero_grad()
        prediction = model(user_id.to(torch.int32), movie_id.to(torch.int32))
        loss = criterion(prediction, rating.to(torch.float32))
        loss.backward()
        optimizer.step()
    # endFor
# endFor


def get_recommendations(model, user_id, top_k=18):
    user_emb = model.user_embeddings(torch.tensor(user_id))
    movie_embs = model.movie_embeddings.weight
    scores = torch.matmul(user_emb, movie_embs.t())
    print("Rear Window")
    print(scores[53])
    print("Cinema Paradiso")
    print(scores[54])
    print("Modern Times")
    print(scores[56])
    print("Witness for the Prosecution")
    print(scores[79])

    print("++++++++++")
    print("The Dark Knight")
    print(scores[3])
    print("Avengers: Endgame")
    print(scores[60])
    print("Avengers: Infinity War")
    print(scores[64])
    print("Spider-Man: Into the Spider-Verse")
    print(scores[69])
    _, top_movie_indices = torch.topk(scores, top_k)
    return top_movie_indices
# endDef


cur = db.cursor()
res = cur.execute("SELECT * FROM User")
users = res.fetchall()

print(get_recommendations(model, 36))
print("==================================================\n\n")
recommendations = get_recommendations(model, 37)
print(recommendations)
recommendations = get_recommendations(model, 46)
print(recommendations)
