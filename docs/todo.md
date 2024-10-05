# Todo

## Todo
### Front-End
- [ ] Wireframes
- [ ] Prototypes
- [ ] Video Player needs a pause/play button to stop the timer
- [ ] Password Recovery Option (clemson email)
- [ ]
### Backend
- [ ]
- [ ]
### Data Analysis
- [ ] Calculate Average Watch Time Percentage metrics for a video
- [ ] Calculate Average Ratings of Each Movie
- [ ] 
- [ ] 
- [ ] 

## In Progress
- [ ] 

## Done
- [x] User sessions are bleeding into each other due to the constant shared context
- [x] Video Player that tracks runtime and regularly updates the WatchEvent table
- [x] Scrape Video Data from IMDB & Rotten Tomatoes (Description, Thumbnail, Runtime, MPA Rating)
### Database Preparation
- [x] Build metadata into a Video table
- [x] User Authentication, (Register, Login, Logout)
- [x] User Table (UserId, Email, PasswordHash, DateAdded, DateModified, FullName)
- [x] Rating Table (VideoId, UserId, Value, DateAdded, DateModified)
- [x] WatchEvent Table (VideoId, UserId, ProgressSeconds, DateAdded, DateModified)
- [x] WatchQueue Table (VideoId, UserId, DateAdded, DateModified)
- [x] 
