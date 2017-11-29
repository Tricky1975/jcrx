/*
  jcrx.go
  
  version: 17.11.29
  Copyright (C) 2017 Jeroen P. Broks
  This software is provided 'as-is', without any express or implied
  warranty.  In no event will the authors be held liable for any damages
  arising from the use of this software.
  Permission is granted to anyone to use this software for any purpose,
  including commercial applications, and to alter it and redistribute it
  freely, subject to the following restrictions:
  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.
  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.
  3. This notice may not be removed or altered from any source distribution.
*/
package main

/*
 
 
   This is only a very small tool, which has been written to act as a 
   dependency for tools to which JCR6 cannot be linked so easily.
   For example LOVE2D.
   
   The first line it returns will always be OK when the operation
   was succesful and OK has not been found here, all output should
   be considered as an error message.
   
   It can basically output the JCR6 file's directory tree in
   both Lua as Python code, and if you have more scripting languages
   which should be supported, lemme know
   
   This tool has been licensed under the terms of the zlib license
   The JCR6 modules have been licensed under the Mozilla Public License
   
*/

import (
	"os"
	"fmt"
	"trickyunits/mkl"
)

func main(){
mkl.Version("jcrx - jcrx.go","17.11.29")
mkl.Lic    ("jcrx - jcrx.go","ZLib License")
	if len(os.Args)<2 {
			fmt.Println("OK")
			fmt.Println(mkl.Newest())
			fmt.Println("Built on sources:")
			fmt.Println(mkl.ListAll())
		}
}
