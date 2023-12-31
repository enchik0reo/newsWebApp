package notifier

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"newsWebApp/app/newsService/internal/models"
	"newsWebApp/app/newsService/internal/services"
	"newsWebApp/app/newsService/internal/storage"
)

type ArticleStorage interface {
	Save(ctx context.Context, article models.Article) error
	NewestNotPosted(ctx context.Context) (*models.Article, error)
	LatestPosted(ctx context.Context, limit int) ([]models.Article, error)
	MarkPosted(ctx context.Context, id int64) (time.Time, error)
}
type Saver interface {
	SaveArticleFromUser(ctx context.Context, userID int64, link string) error
}

type Notifier struct {
	articles ArticleStorage
	saver    Saver

	articlesLimit int
	log           *slog.Logger
}

func New(
	articles ArticleStorage,
	saver Saver,
	articlesLimit int,
	log *slog.Logger,
) *Notifier {
	return &Notifier{
		articles:      articles,
		saver:         saver,
		articlesLimit: articlesLimit,
		log:           log,
	}
}

func (n *Notifier) SaveArticleFromUser(ctx context.Context, userID int64, link string) error {
	const op = "services.notifier.save_article_from_user"

	if err := n.saver.SaveArticleFromUser(ctx, userID, link); err != nil {
		switch {
		case errors.Is(err, services.ErrArticleSkipped):
			n.log.Debug("Can't save article from user", "err", err.Error())
			return services.ErrArticleSkipped
		case errors.Is(err, services.ErrArticleExists):
			n.log.Debug("Can't save article from user", "err", err.Error())
			return services.ErrArticleExists
		default:
			n.log.Error("Can't save article from user", "err", err.Error())
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (n *Notifier) SelectPostedArticles(ctx context.Context) ([]models.Article, error) {
	const op = "services.notifier.select_posted_articles"

	articles, err := n.articles.LatestPosted(ctx, n.articlesLimit)
	if err != nil || len(articles) == 0 {
		switch {
		case len(articles) == 0:
			n.log.Debug("There are no latest posted articles")
			return nil, services.ErrNoPublishedArticles
		case errors.Is(err, storage.ErrNoSources):
			n.log.Debug("Can't get latest posted articles", "err", err.Error())
			return nil, services.ErrNoPublishedArticles
		default:
			n.log.Error("Can't get latest posted articles", "err", err.Error())
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return articles, nil
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) (*models.Article, error) {
	const op = "services.notifier.select_and_send_article"

	article, err := n.articles.NewestNotPosted(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrNoNewArticles) {
			n.log.Debug("Can't get last article", "err", err.Error())
			return nil, services.ErrNoNewArticle
		}
		n.log.Error("Can't get last article", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	postedAt, err := n.articles.MarkPosted(ctx, article.ID)
	if err != nil {
		n.log.Error("Can't mark article as posted", "err", err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	article.PostedAt = postedAt

	return article, nil
}
