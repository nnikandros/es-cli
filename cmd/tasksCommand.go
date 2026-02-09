package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"serde"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/tasks/get"
	"github.com/spf13/cobra"
)

type TasksCmd = *cobra.Command

func tasksCmdFunc(es *elasticsearch.TypedClient) TasksCmd {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "tasks",
		Long:  `A longer description that spans multiple lines and likely contains`,

		RunE: runTasksCmdFunc(es),
		// SilenceUsage:      true,
		// ValidArgsFunction: ValidArgsFuncAutoCompletion(es),

	}

	return cmd

}

func addTasksFlags(tasksCmd TasksCmd) TasksCmd {
	tasksCmd.Flags().BoolP("cancel", "c", false, "this will cancel the running task.")
	tasksCmd.Flags().IntP("watch", "w", 0, "watch mode")

	// tasksCmd.Flags().String("task_id", "", "task idd")
	// tasksCmd.Flags().BoolP("ping", "p", false, "ping flag. Pining the index asserts that you can connect to the index. It makes a match_all query and assets that the response is OK")
	return tasksCmd

}

func TasksCmdFunc(es *elasticsearch.TypedClient) TasksCmd {
	tasks := addTasksFlags(tasksCmdFunc(es))

	return tasks
}

func runTasksCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		r, err := getTask(es, args[0])
		if err != nil {
			return fmt.Errorf("at getTask %w", err)
		}

		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}
}

func getTask(es *elasticsearch.TypedClient, taskid string) (*get.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := es.Tasks.Get(taskid).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("at getting task information %w", err)
	}

	return r, nil

}
