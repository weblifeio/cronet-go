package cronet

// #include <stdlib.h>
// #include <stdbool.h>
// #include <cronet_c.h>
// extern void cronetExecutorExecute(Cronet_ExecutorPtr self,Cronet_RunnablePtr command);
import "C"

import (
	"sync"
	"unsafe"
)

var (
	executorAccess sync.RWMutex
	executors      = make(map[uintptr]ExecutorExecuteFunc)
)

func NewExecutor(executeFunc ExecutorExecuteFunc) Executor {
	ptr := C.Cronet_Executor_CreateWith((*[0]byte)(C.cronetExecutorExecute))
	executorAccess.Lock()
	executors[uintptr(unsafe.Pointer(ptr))] = executeFunc
	executorAccess.Unlock()
	return Executor{ptr}
}

func (e Executor) Destroy() {
	C.Cronet_Executor_Destroy(e.ptr)
	executorAccess.Lock()
	delete(executors, uintptr(unsafe.Pointer(e.ptr)))
	executorAccess.Unlock()
}

//export cronetExecutorExecute
func cronetExecutorExecute(self C.Cronet_ExecutorPtr, command C.Cronet_RunnablePtr) {
	executorAccess.RLock()
	executeFunc := executors[uintptr(unsafe.Pointer(self))]
	executorAccess.RUnlock()
	if executeFunc != nil {
		executeFunc(Executor{self}, Runnable{command})
	}
}
