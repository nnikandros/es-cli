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
	tasksCmd.Flags().IntP("watch", "w", 0, "watch mode, number of seconds to watch")

	return tasksCmd

}

func TasksCmdFunc(es *elasticsearch.TypedClient) TasksCmd {
	tasks := addTasksFlags(tasksCmdFunc(es))

	return tasks
}

func runTasksCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		w, _ := cmd.Flags().GetInt("watch")

		switch w {
		case 0:
			r, err := getTask(es, args[0])
			if err != nil {
				return fmt.Errorf("at getTask %w", err)
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
				return serde.SerializingError(err)
			}

			return nil
		default:

			newEncoder := json.NewEncoder(cmd.OutOrStdout())
			ch := watchTask(es, args[0], w)

			for r := range ch {
				if err := newEncoder.Encode(r); err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "error at encoding response %v", err)
				}

				newEncoder.Encode(r)
			}

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

func watchTask(es *elasticsearch.TypedClient, taskid string, watchFrequency int) <-chan *get.Response {

	resultsChannel := make(chan *get.Response)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)

	ticker := time.Tick(time.Duration(watchFrequency) * time.Second)

	go func() {
		defer cancel()
		defer close(resultsChannel)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				r, err := getTask(es, taskid)
				if err != nil {
					return
				}

				resultsChannel <- r

				if r.Completed {
					return
				}

			}
		}

	}()

	return resultsChannel

}
