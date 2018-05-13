/*
  jcrx.go
  
  version: 18.05.13
  Copyright (C) 2017, 2018 Jeroen P. Broks
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
	"strings"
	"strconv"
	"trickyunits/mkl"
	"trickyunits/qff"
//	"trickyunits/qstr"
	"trickyunits/jcr6/jcr6main"
_	"trickyunits/jcr6/jcr6zlib"
_	"trickyunits/jcr6/jcr6lzw"
_	"trickyunits/jcr6/jcr6flate"
_	"trickyunits/jcr6/jcr6lzma"
_	"trickyunits/jcr6/jcr6realdir"
)

type tsec struct{
	run func() (bool,string)
}
type tdot struct{
	run func(j jcr6main.TJCR6Dir) (bool,string)
}

var section =make(map[string]*tsec)//map[string] tsec
var dirout = make(map[string]*tdot)
var osargs = [] string{}

func defit(){
	// test
	section["test"]=&tsec{}
	section["test"].run = func() (bool,string){
		//fmt.Println("TEST!")
		return true,"TEST"
	}
	// dirout
	section["dirout"]=&tsec{}
	section["dirout"].run = func() (bool,string){
		if len(osargs)<4 {
			return false,"Invalid dirout!"
		}
		if _,ok:=dirout[osargs[3]]; !ok{
			return false,"I don't know how to script out the directory in: "+osargs[3]
		}
		j := jcr6main.Dir(osargs[2])
		if jcr6main.JCR6Error!="" {
			return false,jcr6main.JCR6Error
		}
		rb,rs := dirout[osargs[3]].run(j)
		return rb,rs
	}
	
	// dirout: Lua
	dirout["lua"]=&tdot{}
	dirout["lua"].run = func(j jcr6main.TJCR6Dir) (bool,string){
		ret:="local ret={}\nlocal t={}\n\n"
		//ret = ret + fmt.Sprintf("ret.fat = { size = %d, csize = %s, storage='%s' }\n",j.fatsize,j.fatcsize,j.fatstorage)
		dl:=jcr6main.EntryList(j)
		for i:=0;i<len(dl);i++{
			e:=jcr6main.Entry(j,dl[i])
			if jcr6main.JCR6Error!="" {
				return false,jcr6main.JCR6Error
			}
			if dl[i]!=""{
				dl[i] = strings.Replace(dl[i],"'","\\'",-1)
				ret = ret + fmt.Sprintf("\nt = {} ret['%s']=t -- Entry #%d\n",strings.ToUpper(dl[i]),i+1)
				ret = ret + fmt.Sprintf("t.entry          = '%s'\n",dl[i])
				ret = ret + fmt.Sprintf("t.mainfile       = '%s'\n",strings.Replace(e.Mainfile,"\\","/",-1))
				ret = ret + fmt.Sprintf("t.offset         = %d\n",e.Offset)
				ret = ret + fmt.Sprintf("t.size           = %d\n",e.Size)
				ret = ret + fmt.Sprintf("t.compressedsize = %d\n",e.Compressedsize)
				ret = ret + fmt.Sprintf("t.storage        = '%s'\n",e.Storage)
				ret = ret + fmt.Sprintf("t.author         = %s\n",strconv.QuoteToASCII(e.Author))
				ret = ret + fmt.Sprintf("t.notes          = %s\n",strconv.QuoteToASCII(e.Notes))
				ret = ret +             "t.data           = {}\n"
				for ks, vs := range e.Datastring {
					ret = ret + fmt.Sprintf("\tt.data['%s'] = %s\n",ks,strconv.QuoteToASCII(vs))
				}
				for ki, vi := range e.Dataint {
					ret = ret + fmt.Sprintf("\tt.data['%s'] = %d\n",ki,vi)
				}
				for kb, vb := range e.Databool {
					if vb {
						ret = ret + fmt.Sprintf("\tt.data['%s'] = true\n",kb)
					} else {
						ret = ret + fmt.Sprintf("\tt.data['%s'] = false\n",kb)
					}
				}
			}
		}
		ret = ret + "\n\nreturn ret\n"
		return true,ret
	}
	
	// typeout:
	section["typeout"]=&tsec{}
	section["typeout"].run = func() (bool,string) {
		if len(osargs)<4 { return false,"Invalid typeout" }
		j:=jcr6main.Dir(osargs[2])
		if jcr6main.JCR6Error!="" { return false,jcr6main.JCR6Error }
		b:=jcr6main.JCR_B(j,osargs[3])
		if jcr6main.JCR6Error!="" { return false,jcr6main.JCR6Error }
		return true,string(b)
	}
	
	// transform:
	section["transform"]=&tsec{} // transforms a directory into a JCR file and destroys the original directory if succesful!
	section["transform"].run = func() (bool,string) {
		if len(osargs)<3 { return false,"Hey! What do you want to transform?" }
		origineel:=osargs[2]
		if !qff.IsDir(origineel) { return false,"Original is not a directory, or it doesn't exist at all!" }
		orij:=jcr6main.Dir(osargs[2])
		ret:=""
		tarj:=jcr6main.JCR_Create(osargs[2]+".jcr","BRUTE")
		for ek,ev:=range orij.Entries {
			ret+="Freezing: "+ek+" ... "
			o,c,a:=tarj.AddFile(ev.Mainfile,ev.Entry,"BRUTE","jcrx user","created with jcrx. license set by app using this as dependency")
			ret+=fmt.Sprintf("(%d => %d) %s\n",o,c,a)
		}
		tarj.Close()
		destroy:=true
		if len(osargs)>3 { if osargs[3]=="KEEPME" {destroy=false}}
		if destroy {
			ret+="\nDestroying original: "+origineel+"\n\n"
			os.RemoveAll(origineel)
		}
		return true,ret
	}
	
	// extract:
	section["extract"] = &tsec{}
	section["extract"].run = func() (bool,string) {
		if len(osargs)<5 { return false,"Too little parameters for extraction" }
		jcrfil:=osargs[2]
		source:=osargs[3]
		target:=osargs[4]
		jcr   :=jcr6main.Dir(jcrfil);			if jcr6main.JCR6Error!="" { return false,jcr6main.JCR6Error }
		b     :=jcr6main.JCR_B(jcr,source);		if jcr6main.JCR6Error!="" { return false,jcr6main.JCR6Error }
		bt,err:=os.Create(target);				if err!=nil { return false,err.Error() } ;defer bt.Close()
		bt.Write(b)	;							if err!=nil { return false,err.Error() }
		return true,"Nobody expects the Spanish Inquisition!"
	}
	
	section["getblock"] = &tsec{}
	section["getblock"].run = func() (bool,string) {
		// parsing parameters
		if len(osargs)<7 { return false,"Too little parameters for getblock" }
		offset, err2 := strconv.ParseInt(osargs[2], 10, 64); if err2!=nil {return false,"Invalid offset" } // 2 = offset
		csize , err3 := strconv.ParseInt(osargs[3], 10, 64); if err3!=nil {return false,"Invalid offset" } // 3 = compressed size
		 size , err2 := strconv.ParseInt(osargs[4], 10, 64); if err2!=nil {return false,"Invalid offset" } // 4 = true size
		storage      := osargs[5]                                                                          // 5 = storage method
		mainfile     := osargs[6]                                                                          // 6 = main file
		// declare bank
		pb := make([]byte,csize); 
		// read the compressed bank
		bt,err := os.Open(mainfile)
		defer bt.Close()
		if err!=nil {
			return false,fmt.Sprintf("Error while opening resource file: %s",mainfile)
		}
		bt.Seek(int64(offset),0)
		bt.Read(pb)
		// unpack compressed bank
		var ub []byte
		if stdrv,ok:=jcr6main.JCR6StorageDrivers[storage];ok{
			ub = stdrv.Unpack(pb,int(size))
		} else {
			return false,fmt.Sprintf("Storage algorith %s doesn't exist",storage)
		}
		return true,string(ub)
	}
}


func main(){
mkl.Version("jcrx - jcrx.go","18.05.13")
mkl.Lic    ("jcrx - jcrx.go","ZLib License")
	for _,arg := range os.Args{ osargs = append(osargs,strings.Replace(arg,"DIE_VIEZE_VUILE_FUCK_WINDOWS_HEEFT_EEN_SPATIEBALK_NODIG"," ",-1)) }
	defit()
	if len(osargs)<2 {
		fmt.Println("OK")
		fmt.Println(mkl.Newest())
		fmt.Println("Built on sources:")
		fmt.Println(mkl.ListAll())
	} else { 
		if sec,ok:= section[osargs[1]]; ok{
			success,outdata:=sec.run()
			if success{
				fmt.Println("OK")
			} else {
				fmt.Println("ERROR!")
			}
			fmt.Println(outdata)
		} else {
			fmt.Printf("ERROR!\nI don't know what you mean by %s\n\nDid you spell it right! And please note that I only understand 'lower case' commands!\n",osargs[1])
		}
	}
}
