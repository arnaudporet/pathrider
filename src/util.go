// Copyright (C) 2019 Arnaud Poret
// This work is licensed under the GNU General Public License.
// To view a copy of this license, visit https://www.gnu.org/licenses/gpl.html.

// WARNING The functions in the present file do not fully handle exceptions and
// errors. Instead, they assume that such handling is performed upstream by
// top-level functions of pathrider. Consequently, be careful if using them
// "as is" outside of pathrider.

package main
import (
    "errors"
    "os"
    "strings"
)
func CopyList(list []string) []string {
    var (
        y []string
    )
    y=make([]string,len(list))
    copy(y,list)
    return y
}
func CopyList2(list2 [][]string) [][]string {
    var (
        i int
        y [][]string
    )
    y=make([][]string,len(list2))
    for i=range list2 {
        y[i]=CopyList(list2[i])
    }
    return y
}
func IsInList(list []string,thatElement string) bool {
    var (
        found bool
        element string
    )
    found=false
    for _,element=range list {
        if element==thatElement {
            found=true
            break
        }
    }
    return found
}
func IsInList2(list2 [][]string,thatList []string) bool {
    var (
        found bool
        list []string
    )
    found=false
    for _,list=range list2 {
        if ListEq(list,thatList) {
            found=true
            break
        }
    }
    return found
}
func ListEq(list1,list2 []string) bool {
    var (
        eq bool
        i int
    )
    eq=true
    if len(list1)!=len(list2) {
        eq=false
    } else {
        for i=range list1 {
            if list1[i]!=list2[i] {
                eq=false
                break
            }
        }
    }
    return eq
}
func WriteText(textFile string,text []string) error {
    var (
        err error
        file *os.File
    )
    if len(text)==0 {
        err=errors.New("empty before writing")
    } else {
        file,err=os.Create(textFile)
        defer file.Close()
        if err==nil {
            _,err=file.WriteString(strings.Join(text,"\n")+"\n")
        }
    }
    return err
}
