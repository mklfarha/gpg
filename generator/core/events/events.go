package events

import (
	"context"
	"fmt"
	"os"
	"path"

	"slices"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func GenerateCoreEvents(ctx context.Context, rootPath string, project entity.Project) error {

	if !project.Events.Enabled {
		fmt.Printf("--[GPG][Events] Events disabled skipping\n")
		return nil
	}

	projectDir := generator.ProjectDir(ctx, rootPath, project)
	eventsDir := path.Join(projectDir, generator.EVENTS_REPO_DIR)

	err := os.RemoveAll(eventsDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting core/events directory\n")
	}
	fmt.Printf("--[GPG][Events] Generating core/events module\n")
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(eventsDir, "events.go"),
		TemplateName: path.Join("core", "events", "events"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(eventsDir, "types.go"),
		TemplateName: path.Join("core", "events", "events_types"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(eventsDir, "produce.go"),
		TemplateName: path.Join("core", "events", "events_produce"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(eventsDir, "entity.go"),
		TemplateName: path.Join("core", "events", "events_entity"),
	})
	if err != nil {
		return err
	}

	return nil
}

func ShouldPublishEvents(project entity.Project, identifier string) bool {
	if !project.Events.Enabled {
		return false
	}

	if project.Events.AllEntities {
		return true
	}

	if slices.Contains(project.Events.EntityIdentifiers, identifier) {
		return true
	}

	return false
}
