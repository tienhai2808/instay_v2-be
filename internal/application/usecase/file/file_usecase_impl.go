package usecase

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fileUseCaseImpl struct {
	cfg     config.MinIOConfig
	client  *s3.Client
	pClient *s3.PresignClient
	log     *zap.Logger
}

func NewFileUseCase(
	cfg config.MinIOConfig,
	stor *s3.Client,
	log *zap.Logger,
) FileUseCase {
	pClient := s3.NewPresignClient(stor)
	return &fileUseCaseImpl{
		cfg,
		stor,
		pClient,
		log,
	}
}

func (u *fileUseCaseImpl) CreateUploadURLs(ctx context.Context, req dto.UploadPresignedURLsRequest) ([]*dto.UploadPresignedURLResponse, error) {
	result := make([]*dto.UploadPresignedURLResponse, 0, len(req.Files))

	for _, file := range req.Files {
		name := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)

		key := fmt.Sprintf("%s-%s%s", uuid.NewString(), utils.GenerateSlug(name), ext)
		presignedRes, err := u.pClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(u.cfg.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(file.ContentType),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			u.log.Error("generate upload presigned URL failed", zap.String("content_type", file.ContentType), zap.Error(err))
			return nil, err
		}

		result = append(result, &dto.UploadPresignedURLResponse{
			Key: key,
			Url: presignedRes.URL,
		})
	}

	return result, nil
}

func (u *fileUseCaseImpl) CreateViewURLs(ctx context.Context, req dto.ViewPresignedURLsRequest) ([]*dto.ViewPresignedURLResponse, error) {
	result := make([]*dto.ViewPresignedURLResponse, 0, len(req.Keys))

	for _, key := range req.Keys {
		if _, err := u.client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(u.cfg.Bucket),
			Key:    aws.String(key),
		}); err != nil {
			var keyNotFound *types.NotFound
			if errors.As(err, &keyNotFound) {
				result = append(result, nil)
				continue
			}
			u.log.Error("file check failed", zap.Error(err))
			return nil, err
		}

		presignedReq, err := u.pClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(u.cfg.Bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = 15 * time.Minute
		})
		if err != nil {
			u.log.Error("generate view presigned URL failed", zap.Error(err))
			return nil, err
		}

		result = append(result, &dto.ViewPresignedURLResponse{
			Url: presignedReq.URL,
		})
	}

	return result, nil
}
