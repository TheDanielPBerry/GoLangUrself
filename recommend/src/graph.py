import matplotlib.pyplot as plt
import numpy as np
import sqlite3

plt.style.use('_mpl-gallery')



db = sqlite3.connect('../../db/westflix.db')
cur = db.cursor()
res = cur.execute("SELECT * FROM AverageWatchTimes ORDER BY TotalWatchPercentage LIMIT 30")
watchTimes = res.fetchall()

titles = []
totalPercentages = []
for movie in watchTimes:
    titles.append(movie[1])
    totalPercentages.append(movie[3])
# endFor

titles.append('')
totalPercentages.append(0)

cur = db.cursor()
res = cur.execute("SELECT * FROM AverageWatchTimes ORDER BY TotalWatchPercentage DESC LIMIT 30")
watchTimes = res.fetchall()

for movie in reversed(watchTimes):
    titles.append(movie[1])
    totalPercentages.append(movie[3])
# endFor

# plot
fig, ax = plt.subplots()

ax.barh(titles, totalPercentages, height=1, edgecolor="white", linewidth=0.2)

ax.set(xlim=(0, 23), xticks=np.arange(1, 23),
       ylim=(0, 61), yticks=titles)
ax.set_xlabel('Total WatchTime Percentages')
ax.set_ylabel('Movie Title')

fig.xmargin = 2;
plt.show()

