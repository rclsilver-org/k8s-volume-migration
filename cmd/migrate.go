package cmd

import (
	"os"

	"github.com/rclsilver-org/k8s-volume-migration/pkg/directories"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	migrateSrcDirFlag      = "source-directory"
	migrateSrcDirShortFlag = "s"
	migrateSrcDirDefault   = ""

	migrateDstDirFlag      = "destination-directory"
	migrateDstDirShortFlag = "d"
	migrateDstDirDefault   = ""

	migrateChgOwnerFlag      = "owner"
	migrateChgOwnerShortFlag = "u"
	migrateChgOwnerDefault   = ""

	migrateChgGroupFlag      = "group"
	migrateChgGroupShortFlag = "g"
	migrateChgGroupDefault   = ""
)

var (
	migrateSrcDir   string
	migrateDstDir   string
	migrateChgOwner string
	migrateChgGroup string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		srcIsEmpty, err := directories.IsEmpty(ctx, migrateSrcDir)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Fatal("unable to read the content of the source directory")
		}
		if srcIsEmpty {
			logrus.WithContext(ctx).Fatal("the source directory is empty")
		}

		dstIsEmpty, err := directories.IsEmpty(ctx, migrateDstDir)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Fatal("unable to read the content of the destination directory")
		}
		if !dstIsEmpty {
			logrus.WithContext(ctx).Info("the destination directory is not empty, migration is not required")
			os.Exit(0)
		}

		logrus.WithContext(ctx).Infof("copying data from %q to %q", migrateSrcDir, migrateDstDir)

		if err := directories.CopyDir(ctx, migrateSrcDir, migrateDstDir, migrateChgOwner, migrateChgGroup); err != nil {
			logrus.WithContext(ctx).WithError(err).Fatal("error while copying data")
		}

		logrus.WithContext(ctx).Infof("all the data have been migrated from %q to %q", migrateSrcDir, migrateDstDir)
	},
}

func init() {
	migrateCmd.PersistentFlags().StringVarP(&migrateSrcDir, migrateSrcDirFlag, migrateSrcDirShortFlag, migrateSrcDirDefault, "Source data directory")
	migrateCmd.MarkPersistentFlagRequired(migrateSrcDirFlag)

	migrateCmd.PersistentFlags().StringVarP(&migrateDstDir, migrateDstDirFlag, migrateDstDirShortFlag, migrateDstDirDefault, "Destination data directory")
	migrateCmd.MarkPersistentFlagRequired(migrateDstDirFlag)

	migrateCmd.PersistentFlags().StringVarP(&migrateChgOwner, migrateChgOwnerFlag, migrateChgOwnerShortFlag, migrateChgOwnerDefault, "Change files owner")
	migrateCmd.PersistentFlags().StringVarP(&migrateChgGroup, migrateChgGroupFlag, migrateChgGroupShortFlag, migrateChgGroupDefault, "Change files group")

	rootCmd.AddCommand(migrateCmd)
}
