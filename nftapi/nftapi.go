package nftapi

/*
#cgo CFLAGS: -I/home/ap/nftproxy/include
#cgo LDFLAGS: -L/home/ap/nftproxy/lib -lfxnapi
#include "fxapi.h"
#include "fterror.h"
#include "stdio.h"
#include "stdlib.h"

*/
import "C"
import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"nftserver/protocol"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Ntransfer(req *protocol.NFTParamReq) protocol.NFTParamRet {

	ret := protocol.NFTParamRet{}

	//cstr_filepath := C.CString(fmt.Sprintf("%s/%s", os.Getenv("NFTFILEPATH"), req.FilePath))

	cstr_routemode := C.CString(req.RouteMode)
	cstr_routerid := C.CString(req.RouteID)
	cstr_direction := C.CString(req.Direction)
	cstr_apptype := C.CString(req.AppType)
	cstr_filemode := C.CString(req.FileMode)
	cstr_usermsg := C.CString(req.UserMsg)
	tsfid := make([]byte, 40)
	cstr_user := C.CString(req.User)
	cstr_password := C.CString(req.Password)

	defer C.free(unsafe.Pointer(cstr_routemode))
	defer C.free(unsafe.Pointer(cstr_routerid))
	defer C.free(unsafe.Pointer(cstr_direction))
	defer C.free(unsafe.Pointer(cstr_apptype))
	defer C.free(unsafe.Pointer(cstr_filemode))
	defer C.free(unsafe.Pointer(cstr_usermsg))
	defer C.free(unsafe.Pointer(cstr_user))
	defer C.free(unsafe.Pointer(cstr_password))
	/*  修改lst文件，S/R 两种情况
	 *  S 修改第一列的路径信息
	 *  R 修改第二列的路径信息
	 */
	e, msg := modifyListFile(req.FilePath, req.Direction)
	if e != nil {
		ret.RetCode = -1
		ret.RetMsg = msg
		ret.TransferID = ""
		ret.TransferState = ""
		return ret
	}
	cstr_filepath := C.CString(msg)

	defer C.free(unsafe.Pointer(cstr_filepath))
	defer func() {
		if r := recover(); r != nil {
			log.Println("meet panic: ", r)
		}
	}()
	r := C.ntransfer(cstr_filepath, cstr_routemode, cstr_routerid, cstr_direction,
		cstr_apptype, cstr_filemode, cstr_usermsg,
		(*C.char)(unsafe.Pointer(&tsfid[0])),
		cstr_user, cstr_password)
	if r != 0 {
		ret.RetCode = int32(r)
		ret.RetMsg = "投递失败"
		ret.TransferID = ""
		ret.TransferState = ""
	} else {
		ret.RetMsg = fmt.Sprintf("任务投递成功,taskid=[%s]", string(tsfid))
		ret.RetCode = int32(r)
		ret.TransferID = string(tsfid)
		ret.TransferState = ""
	}

	return ret
}

func tempfile(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return prefix + strconv.Itoa(rand.Intn(1000))
}

func modifyListFile(listfile string, flag string) (error, string) {
	sfile := filepath.Join(os.Getenv("NFTFILEPATH"), listfile)
	nfile := filepath.Join(os.Getenv("NFTFILEPATH"), tempfile(listfile+"-"))

	fd, err := os.Open(sfile)
	if err != nil {
		return err, "打开文件失败" + listfile
	}
	defer fd.Close()

	fd1, err := os.Create(nfile)
	if err != nil {
		return err, "创建文件失败" + nfile
	}
	defer fd1.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, " ")

		a := 0
		if flag == "R" && len(arr) > 1 {
			a = 1
		}
		arr[a] = filepath.Join(os.Getenv("NFTFILEPATH"), filepath.Base(arr[a]))
		if _, err := fd1.WriteString(strings.Join(arr, " ")); err != nil {
			return err, "写入新list文件出错"
		}
	}
	return nil, nfile
}

func GetTsfResponse(req *protocol.NFTParamReq) protocol.NFTParamRet {
	ret := protocol.NFTParamRet{}

	cstr_transferid := C.CString(req.TransferID)
	info := make([]byte, 1024)

	defer C.free(unsafe.Pointer(cstr_transferid))

	r := C.getTsfResponse(cstr_transferid, (*C.char)(unsafe.Pointer(&info[0])))
	if r != 0 {
		ret.RetCode = int32(r)
		ret.RetMsg = "获取任务信息出错!"
		ret.TransferID = req.TransferID
		ret.TransferState = ""
	} else {
		ret.RetMsg = string(info)
		ret.RetCode = int32(r)
		ret.TransferID = req.TransferID
		ret.TransferState = ""
	}

	return ret
}

func GetTsfState(req *protocol.NFTParamReq) protocol.NFTParamRet {
	ret := protocol.NFTParamRet{}

	cstr_transferid := C.CString(req.TransferID)
	state := make([]byte, 1024)

	defer C.free(unsafe.Pointer(cstr_transferid))

	r := C.getTsfState(cstr_transferid, (*C.char)(unsafe.Pointer(&state[0])))
	if r != 0 {
		ret.RetCode = int32(r)
		ret.RetMsg = "获取任务信息出错"
		ret.TransferID = req.TransferID
		ret.TransferState = ""
	} else {
		ret.RetMsg = "成功"
		ret.RetCode = int32(r)
		ret.TransferID = req.TransferID
		ret.TransferState = string(state)
	}

	return ret
}

func NCancelTransfer(req *protocol.NFTParamReq) protocol.NFTParamRet {
	ret := protocol.NFTParamRet{}

	cstr_transferid := C.CString(req.TransferID)
	defer C.free(unsafe.Pointer(cstr_transferid))

	r := C.ncancelTransfer(cstr_transferid)
	if r != 0 {
		ret.RetCode = int32(r)
		ret.RetMsg = "取消任务失败"
		ret.TransferID = req.TransferID
		ret.TransferState = ""
	} else {
		ret.RetMsg = "成功"
		ret.RetCode = int32(r)
		ret.TransferID = req.TransferID
		ret.TransferState = ""
	}

	return ret
}
