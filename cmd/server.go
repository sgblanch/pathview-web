/*
Copyright © 2022 Steven Blanchard <sgblanch@uncc.edu>

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serve(cmd *cobra.Command, args []string) {
	log.SetPrefix("[kegg:server]")

	viper.BindEnv("csrf-key")
	viper.SetDefault("listen", "localhost:8000")
	viper.BindEnv("redis.address", "REDIS_ADDR")
	viper.SetDefault("redis.address", "localhost:6379")
	viper.BindEnv("session.auth-key")
	viper.BindEnv("session.enc-key")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := new(server.Router)
	srv := &http.Server{
		Addr:    config.Get().Listen,
		Handler: router.Router(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Print("shutting down..")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("aborted", err)
	}
}
