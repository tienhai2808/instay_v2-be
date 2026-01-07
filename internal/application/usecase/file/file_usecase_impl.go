package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/InstayPMS/backend/internal/application/dto"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type fileUseCaseImpl struct {
	cfg  *config.Config
	stor *minio.Client
	log  *zap.Logger
}

func NewFileUseCase(
	cfg *config.Config,
	str *minio.Client,
	log *zap.Logger,
) FileUseCase {
	return &fileUseCaseImpl{
		cfg,
		str,
		log,
	}
}

func (s *fileUseCaseImpl) CreateUploadURLs(ctx context.Context, req dto.UploadPresignedURLsRequest) ([]*dto.UploadPresignedURLResponse, error) {
	result := make([]*dto.UploadPresignedURLResponse, 0, len(req.Files))

	for _, file := range req.Files {
		name := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)

		key := fmt.Sprintf("%s-%s%s", uuid.NewString(), utils.GenerateSlug(name), ext)
		expiresIn := 15 * time.Minute

		presignedURL, err := s.stor.PresignedPutObject(ctx, s.cfg.MinIO.Bucket, key, expiresIn)
		if err != nil {
			s.log.Error("generate upload presigned URL failed", zap.String("content_type", file.ContentType), zap.Error(err))
			return nil, err
		}

		result = append(result, &dto.UploadPresignedURLResponse{
			Key: key,
			Url: presignedURL.String(),
		})
	}

	return result, nil
}

func (s *fileUseCaseImpl) CreateViewURLs(ctx context.Context, req dto.ViewPresignedURLsRequest) ([]*dto.ViewPresignedURLResponse, error) {
	result := make([]*dto.ViewPresignedURLResponse, 0, len(req.Keys))

	for _, key := range req.Keys {
		if _, err := s.stor.StatObject(ctx, s.cfg.MinIO.Bucket, key, minio.StatObjectOptions{}); err != nil {
			errResponse := minio.ToErrorResponse(err)
			if errResponse.Code == "NoSuchKey" {
				result = append(result, nil)
				continue
			}
			s.log.Error("file check failed", zap.Error(err))
			return nil, err
		}

		expiresIn := 15 * time.Minute
		presignedURL, err := s.stor.PresignedGetObject(ctx, s.cfg.MinIO.Bucket, key, expiresIn, nil)
		if err != nil {
			s.log.Error("generate view presigned URL failed", zap.Error(err))
			return nil, err
		}

		result = append(result, &dto.ViewPresignedURLResponse{
			Url: presignedURL.String(),
		})
	}

	return result, nil
}
