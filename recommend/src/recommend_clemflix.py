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

for epoch in range(100):
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
    _, top_movie_indices = torch.topk(scores, top_k)
    return top_movie_indices
# endDef


cur = db.cursor()
res = cur.execute("SELECT * FROM User")
users = res.fetchall()

for u in users:
    recommendations = get_recommendations(model, u[0])
    print(recommendations)
# endFor
