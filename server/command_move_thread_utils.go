package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
)

// WranglerPostList contains a list of posts and metadata about those posts
type WranglerPostList struct {
	Posts                []*model.Post
	ThreadUserIDs        []string
	EarlistPostTimeStamp int64
	LatestPostTimeStamp  int64
	ContainsAttachments  bool
}

// NumPosts returns the number of posts.
func (wpl *WranglerPostList) NumPosts() int {
	return len(wpl.Posts)
}

func sortedPostsFromPostList(postList *model.PostList) *WranglerPostList {
	wpl := &WranglerPostList{}

	// A separate ID key map to ensure no duplicates.
	idKeys := make(map[string]bool)

	postList.UniqueOrder()
	postList.SortByCreateAt()
	posts := postList.ToSlice()

	if len(posts) == 0 {
		// Something was sorted wrong or an empty PostList was provided.
		return wpl
	}

	for i := range posts {
		p := posts[len(posts)-i-1]

		// Add UserID to metadata if it's new.
		if _, ok := idKeys[p.UserId]; !ok {
			idKeys[p.UserId] = true
			wpl.ThreadUserIDs = append(wpl.ThreadUserIDs, p.UserId)
		}

		// Mark postlist as containing attachments if post has attachment(s).
		if !wpl.ContainsAttachments && len(p.Attachments()) != 0 {
			wpl.ContainsAttachments = true
		}

		wpl.Posts = append(wpl.Posts, p)
	}

	// Set metadata for earliest and latest posts
	wpl.EarlistPostTimeStamp = wpl.Posts[0].CreateAt
	wpl.LatestPostTimeStamp = wpl.Posts[wpl.NumPosts()-1].CreateAt

	return wpl
}

func cleanPost(post *model.Post) {
	post.Id = ""
	post.CreateAt = 0
	post.UpdateAt = 0
	post.EditAt = 0
}
