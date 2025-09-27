/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var port int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the ToDo Backend API",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
		fmt.Println("Server is running on port ", port)
	},
}

func init() {
	runCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number to listen on")
	rootCmd.AddCommand(runCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
