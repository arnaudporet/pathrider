// Copyright (C) 2019 Arnaud Poret
// This work is licensed under the GNU General Public License.
// To view a copy of this license, visit https://www.gnu.org/licenses/gpl.html.
package main
import (
    "flag"
    "fmt"
    "os"
    "strings"
)
func main() {
    var (
        err error
        help,usage,license bool
        command string
        flagSet *flag.FlagSet
    )
    flagSet=flag.NewFlagSet("",flag.ContinueOnError)
    flagSet.Usage=func() {}
    flagSet.BoolVar(&help,"help",false,"")
    flagSet.BoolVar(&help,"h",false,"")
    flagSet.BoolVar(&usage,"usage",false,"")
    flagSet.BoolVar(&usage,"u",false,"")
    flagSet.BoolVar(&license,"license",false,"")
    flagSet.BoolVar(&license,"l",false,"")
    err=flagSet.Parse(os.Args[1:])
    if err!=nil {
        fmt.Println("Error: pathrider: "+err.Error())
    } else if help {
        fmt.Println(strings.Join([]string{
            "",
            "pathrider is a tool for finding paths of interest in networks.",
            "",
            "pathrider currently provides 2 commands:",
            "    * connect: find the paths connecting some nodes of interest in a network",
            "    * stream: find the upstream/downstream paths starting from some nodes of",
            "              interest in a network",
            "",
            "Usage:",
            "    * pathrider [options]",
            "    * pathrider <command> [options] <arguments>",
            "",
            "Positional argument:",
            "    * <command>: connect, stream",
            "",
            "Options:",
            "    * -l/-license: print the GNU General Public License under which pathrider is",
            "    * -u/-usage: print usage only",
            "    * -h/-help: print help",
            "",
            "For command-specific help, run \"pathrider <command> -help\".",
            "",
            "For more information, see https://github.com/arnaudporet/pathrider.",
            "",
        },"\n"))
    } else if usage {
        fmt.Println(strings.Join([]string{
            "",
            "Usage:",
            "    * pathrider [options]",
            "    * pathrider <command> [options] <arguments>",
            "",
            "Positional argument:",
            "    * <command>: connect, stream",
            "",
            "Options:",
            "    * -l/-license: print the GNU General Public License under which pathrider is",
            "    * -u/-usage: print usage only",
            "    * -h/-help: print help",
            "",
        },"\n"))
    } else if license {
        fmt.Println(strings.Join([]string{
            "",
            "pathrider: a tool for finding paths of interest in networks.",
            "Copyright (C) 2019 Arnaud Poret",
            "",
            "This program is free software: you can redistribute it and/or modify it under",
            "the terms of the GNU General Public License as published by the Free Software",
            "Foundation, either version 3 of the License, or (at your option) any later",
            "version.",
            "",
            "This program is distributed in the hope that it will be useful, but WITHOUT ANY",
            "WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A",
            "PARTICULAR PURPOSE. See the GNU General Public License for more details.",
            "",
            "You should have received a copy of the GNU General Public License along with",
            "this program. If not, see <https://www.gnu.org/licenses/>.",
            "",
        },"\n"))
    } else if len(flagSet.Args())==0 {
        fmt.Println("Error: pathrider: missing command, expecting one of: connect, stream")
    } else {
        command=flagSet.Arg(0)
        if command=="connect" {
            Connect()
        } else if command=="stream" {
            Stream()
        } else {
            fmt.Println("Error: pathrider: "+command+": unknown command, expecting one of: connect, stream")
        }
    }
}
