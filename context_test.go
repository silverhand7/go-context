package go_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	backgorund := context.Background()
	fmt.Println(backgorund)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B") // membuat child context
	contextC := context.WithValue(contextA, "c", "C") // membuat child context

	contextD := context.WithValue(contextB, "d", "D") // membuat child context
	contextE := context.WithValue(contextB, "e", "E") // membuat child context

	contextF := context.WithValue(contextC, "f", "F") // membuat child context
	contextG := context.WithValue(contextF, "g", "G") // membuat child context

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f")) // dapat
	fmt.Println(contextF.Value("c")) // dapat, from parent
	fmt.Println(contextF.Value("d")) // tidak dapat, beda parent (nil)
	fmt.Println(contextF.Value("g")) // tidak bisa mengambil data child
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		// try to make infinity loop
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Current Goroutine", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}
	cancel() // mengirim sinyal cancel ke context

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func CreateCounterSlow(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		// try to make infinity loop
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulasi slow process
			}
		}
	}()

	return destination
}

func TestContextWithTimeout(t *testing.T) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel() // jika fungsi selesai dieksekusi sebelum timeout, maka cancel() akan dijalankan

	destination := CreateCounterSlow(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second) // ini cuma untuk memastikan goroutine sudah dimatikan

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	parent := context.Background()
	// jika sudah mencapai waktu deadline, maka langsung dihentikan.
	// deadline dijalankan dengan waktu yang sudah ditentukan
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel() // jika fungsi selesai dieksekusi sebelum timeout, maka cancel() akan dijalankan

	destination := CreateCounterSlow(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second) // ini cuma untuk memastikan goroutine sudah dimatikan

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
