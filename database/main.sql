-- Create the Database
CREATE DATABASE IF NOT EXISTS twitter_clone;

-- Use the created database
USE twitter_clone;

-- Users Table: Stores user information
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Tweets Table: Stores tweets
CREATE TABLE IF NOT EXISTS tweets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Followers Table: Stores user-following relationships
CREATE TABLE IF NOT EXISTS followers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    follower_id INT NOT NULL,
    following_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT UNIQUE(follower_id, following_id) -- A user can follow another user only once
);

-- Likes Table: Stores tweet likes
CREATE TABLE IF NOT EXISTS likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    tweet_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE,
    CONSTRAINT UNIQUE(user_id, tweet_id) -- A user can like a tweet only once
);

-- Retweets Table: Stores tweet retweets
CREATE TABLE IF NOT EXISTS retweets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    tweet_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (tweet_id) REFERENCES tweets(id) ON DELETE CASCADE,
    CONSTRAINT UNIQUE(user_id, tweet_id) -- A user can retweet a tweet only once
);

-- Password Resets Table: Stores password reset tokens for user recovery
CREATE TABLE IF NOT EXISTS password_resets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance (optional but recommended)
CREATE INDEX idx_user_email ON users (email);
CREATE INDEX idx_user_username ON users (username);
CREATE INDEX idx_tweet_user_id ON tweets (user_id);
CREATE INDEX idx_followers_user_id ON followers (follower_id);
CREATE INDEX idx_likes_user_id ON likes (user_id);
CREATE INDEX idx_retweets_user_id ON retweets (user_id);

