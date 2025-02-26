package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gookit/color"
	"github.com/goravel/framework/http"
	"github.com/goravel/framework/support/str"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/filesystem"
)

/*
 * S3 OSS
 * Document: https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/s3
 * More: https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/s3/#putobjectinput-body-field-ioreadseeker-vs-ioreader
 */

type S3 struct {
	ctx      context.Context
	config   config.Config
	disk     string
	doCdnUrl string
	instance *s3.Client
	bucket   string
	url      string
}

func NewS3(ctx context.Context, config config.Config, disk string) (*S3, error) {
	accessKeyId := config.GetString(fmt.Sprintf("filesystems.disks.%s.key", disk))
	accessKeySecret := config.GetString(fmt.Sprintf("filesystems.disks.%s.secret", disk))
	region := config.GetString(fmt.Sprintf("filesystems.disks.%s.region", disk))
	bucket := config.GetString(fmt.Sprintf("filesystems.disks.%s.bucket", disk))
	url := config.GetString(fmt.Sprintf("filesystems.disks.%s.url", disk))
	token := config.GetString(fmt.Sprintf("filesystems.disks.%s.token", disk), "")
	endpoint := config.GetString(fmt.Sprintf("filesystems.disks.%s.endpoint", disk))
	use_path_style_path := fmt.Sprintf("filesystems.disks.%s.use_path_style", disk)
	use_path_style := config.Get(use_path_style_path)
	use_path_style_bool := config.GetBool(use_path_style_path)
	doCdnUrl := config.GetString(fmt.Sprintf("filesystems.disks.%s.do_cdn_url", disk))

	if accessKeyId == "" || accessKeySecret == "" || region == "" || bucket == "" || url == "" {
		return nil, fmt.Errorf("please set %s configuration first", disk)
	}

	options := s3.Options{
		Region: region,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, token)),
	}

	if endpoint != "" {
		options.BaseEndpoint = aws.String(endpoint)
	}

	if use_path_style != nil {
		options.UsePathStyle = use_path_style_bool
	}

	client := s3.New(options)

	return &S3{
		ctx:      ctx,
		config:   config,
		disk:     disk,
		doCdnUrl: doCdnUrl,
		instance: client,
		bucket:   bucket,
		url:      url,
	}, nil
}

func (r *S3) AllDirectories(path string) ([]string, error) {
	var directories []string
	validPath := validPath(path)
	listObjsResponse, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(r.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(validPath),
	})
	if err != nil {
		return nil, err
	}

	for _, commonPrefix := range listObjsResponse.CommonPrefixes {
		prefix := *commonPrefix.Prefix
		directories = append(directories, strings.ReplaceAll(prefix, validPath, ""))

		subDirectories, err := r.AllDirectories(*commonPrefix.Prefix)
		if err != nil {
			return nil, err
		}
		for _, subDirectory := range subDirectories {
			directories = append(directories, strings.ReplaceAll(prefix+subDirectory, validPath, ""))
		}
	}

	return directories, nil
}

func (r *S3) AllFiles(path string) ([]string, error) {
	var files []string
	validPath := validPath(path)
	listObjsResponse, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(validPath),
	})
	if err != nil {
		return nil, err
	}
	for _, object := range listObjsResponse.Contents {
		file := *object.Key
		if !strings.HasSuffix(file, "/") {
			files = append(files, strings.ReplaceAll(file, validPath, ""))
		}
	}

	return files, nil
}

func (r *S3) Copy(originFile, targetFile string) error {
	_, err := r.instance.CopyObject(r.ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(r.bucket),
		CopySource: aws.String(r.bucket + "/" + originFile),
		Key:        aws.String(targetFile),
	})

	return err
}

func (r *S3) Delete(files ...string) error {
	var objectIdentifiers []types.ObjectIdentifier
	for _, file := range files {
		objectIdentifiers = append(objectIdentifiers, types.ObjectIdentifier{
			Key: aws.String(file),
		})
	}

	_, err := r.instance.DeleteObjects(r.ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(r.bucket),
		Delete: &types.Delete{
			Objects: objectIdentifiers,
			Quiet:   aws.Bool(true),
		},
	})

	return err
}

func (r *S3) DeleteDirectory(directory string) error {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}

	listObjectsV2Response, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(directory),
	})
	if err != nil {
		return err
	}
	if len(listObjectsV2Response.Contents) == 0 {
		return nil
	}

	for {
		for _, item := range listObjectsV2Response.Contents {
			_, err = r.instance.DeleteObject(r.ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(r.bucket),
				Key:    item.Key,
			})
			if err != nil {
				return err
			}
		}

		if *listObjectsV2Response.IsTruncated {
			listObjectsV2Response, err = r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
				Bucket:            aws.String(r.bucket),
				ContinuationToken: listObjectsV2Response.ContinuationToken,
			})
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

func (r *S3) Directories(path string) ([]string, error) {
	var directories []string
	validPath := validPath(path)
	listObjsResponse, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(r.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(validPath),
	})
	if err != nil {
		return nil, err
	}
	for _, commonPrefix := range listObjsResponse.CommonPrefixes {
		directories = append(directories, strings.ReplaceAll(*commonPrefix.Prefix, validPath, ""))
	}

	return directories, nil
}

func (r *S3) Exists(file string) bool {
	_, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})

	if err != nil {
		log.Println("error while checking file existance:", err)
		return false
	}

	return err == nil
}

func (r *S3) Files(path string) ([]string, error) {
	var files []string
	validPath := validPath(path)
	listObjsResponse, err := r.instance.ListObjectsV2(r.ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(r.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(validPath),
	})
	if err != nil {
		return nil, err
	}
	for _, object := range listObjsResponse.Contents {
		file := strings.ReplaceAll(*object.Key, validPath, "")
		if file == "" {
			continue
		}

		files = append(files, file)
	}

	return files, nil
}

func (r *S3) Get(file string) (string, error) {
	data, err := r.GetBytes(file)

	return string(data), err
}

func (r *S3) GetBytes(file string) ([]byte, error) {
	resp, err := r.instance.GetObject(r.ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *S3) LastModified(file string) (time.Time, error) {
	resp, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return time.Time{}, err
	}

	l, err := time.LoadLocation(r.config.GetString("app.timezone"))
	if err != nil {
		return time.Time{}, err
	}

	return aws.ToTime(resp.LastModified).In(l), nil
}

func (r *S3) MakeDirectory(directory string) error {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}

	return r.Put(directory, "")
}

func (r *S3) MimeType(file string) (string, error) {
	resp, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(resp.ContentType), nil
}

func (r *S3) Missing(file string) bool {
	return !r.Exists(file)
}

func (r *S3) Move(oldFile, newFile string) error {
	if err := r.Copy(oldFile, newFile); err != nil {
		return err
	}

	return r.Delete(oldFile)
}

func (r *S3) Path(file string) string {
	return file
}

func (r *S3) Put(file string, content string) error {
	// If the file is created in a folder directly, we can't check if the folder exists.
	// So we need to create the folders first.
	if !strings.HasSuffix(file, "/") {
		folders := strings.Split(file, "/")
		for i := 1; i < len(folders); i++ {
			folder := strings.Join(folders[:i], "/")
			if err := r.MakeDirectory(folder); err != nil {
				return err
			}
		}
	}

	mtype := mimetype.Detect([]byte(content))
	_, err := r.instance.PutObject(r.ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(file),
		Body:          strings.NewReader(content),
		ContentLength: aws.Int64(int64(len(content))),
		ContentType:   aws.String(mtype.String()),
		ACL:           types.ObjectCannedACLPublicRead,
	})

	return err
}

func (r *S3) PutFile(filePath string, source filesystem.File) (string, error) {
	return r.PutFileAs(filePath, source, str.Random(40))
}

func (r *S3) PutFileAs(filePath string, source filesystem.File, name string) (string, error) {
	fullPath, err := fullPathOfFile(filePath, source, name)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(source.File())
	if err != nil {
		return "", err
	}

	if err := r.Put(fullPath, string(data)); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (r *S3) Size(file string) (int64, error) {
	resp, err := r.instance.HeadObject(r.ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	})
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func (r *S3) TemporaryUrl(file string, t time.Time) (string, error) {
	presignClient := s3.NewPresignClient(r.instance)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(file),
	}
	presignDuration := func(po *s3.PresignOptions) {
		po.Expires = time.Until(t)
	}
	presignResult, err := presignClient.PresignGetObject(r.ctx, presignParams, presignDuration)
	if err != nil {
		return "", err
	}

	return presignResult.URL, nil
}

func (r *S3) WithContext(ctx context.Context) filesystem.Driver {
	if httpCtx, ok := ctx.(http.Context); ok {
		ctx = httpCtx.Context()
	}

	driver, err := NewS3(ctx, r.config, r.disk)
	if err != nil {
		color.Redf("[S3] init disk error: %+v\n", err)

		return nil
	}

	return driver
}

func (r *S3) Url(file string) string {
	var u string

	if r.doCdnUrl != "" {
		u = strings.TrimSuffix(r.doCdnUrl, "/") + "/" + strings.TrimPrefix(file, "/")
	} else {
		u = strings.TrimSuffix(r.url, "/") + "/" + strings.TrimPrefix(file, "/")
	}

	return u
}
