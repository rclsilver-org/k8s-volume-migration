package directories

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

func IsEmpty(ctx context.Context, path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, fmt.Errorf("unable to read the destination info: %w", err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("the destination is not a directory")
	}

	dir, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("unable to open the destination directory: %w", err)
	}
	defer dir.Close()

	entries, err := dir.Readdirnames(1)
	if err != nil {
		return false, fmt.Errorf("unable to read the destination directory: %w", err)
	}

	return len(entries) == 0, nil
}

func CopyDir(ctx context.Context, srcPath, dstPath, toOwner, toGroup string) error {
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dstPath, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return fmt.Errorf("unable to walk through the directories: %w", err)
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return fmt.Errorf("unable to build the path: %w", err)
		}
		dstFullPath := filepath.Join(dstPath, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(dstFullPath, info.Mode()); err != nil {
				return fmt.Errorf("unable to create the destination directories: %w", err)
			}
		} else {
			logrus.WithContext(ctx).Debugf("copying %q to %q", path, dstFullPath)

			if err := copyFile(ctx, path, dstFullPath); err != nil {
				return fmt.Errorf("error while copying the file %q", path)
			}

			logrus.WithContext(ctx).Infof("copied %q to %q", path, dstFullPath)
		}

		if err := os.Chmod(dstFullPath, info.Mode()); err != nil {
			return err
		}

		if toOwner != "" || toGroup != "" {
			uid, gid := -1, -1

			if toOwner != "" {
				uid, err = strconv.Atoi(toOwner)
				if err != nil {
					return fmt.Errorf("invalid owner uid: %w", err)
				}
			}

			if toGroup != "" {
				gid, err = strconv.Atoi(toGroup)
				if err != nil {
					return fmt.Errorf("invalid group gid: %w", err)
				}
			}

			if err := os.Chown(dstFullPath, uid, gid); err != nil {
				return fmt.Errorf("failed to change ownership: %w", err)
			}
		}

		return nil
	})

	return err
}

const bufferSize = 1024 * 32

func copyFile(ctx context.Context, src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, readErr := srcFile.Read(buf)
		if n > 0 {
			if _, writeErr := dstFile.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	if err := dstFile.Sync(); err != nil {
		return err
	}

	return nil
}
