# Todo

## Todo
### Front-End
- [ ] Wireframes
- [ ] Prototypes
- [ ] Video Player needs a pause/play button to stop the timer
- [ ] Password Recovery Option (clemson email)
- [x] Video List View
- [x] Navbar
- [x] Fancy video row scrolling
- [ ] Settings Page
- [ ] Account Dropdown

#### Display Following Categories on Home Page in a Row:
- [x] Most Popular
- [x] Continue Watching
- [ ] For You
- [ ] Queue List
- [ ] Format Preview Tray
- [ ] Close Preview tray on click
- [ ] Animate Preview Tray

### Backend
- [ ] Todo Like/Dislike System
- [ ] Add to Queue
- [ ]
### Data Analysis
- [x] Calculate Average Watch Time Percentage metrics for a video
- [x] Calculate Total Watch Time for each video
- [ ] Calculate Average Ratings of Each Movie
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
