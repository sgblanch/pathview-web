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
	"log"

	"github.com/sgblanch/pathview-web/internal/kegg"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetPrefix("[kegg:update] ")
		kegg.Update()
	},
}

func init() {
	keggCmd.AddCommand(updateCmd)
}
