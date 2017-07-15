# Beanstalkd RPC

## Description
Beanstalkd RPC for Go with gin-gonic syntax

## Installation
    go get github.com/MarinX/beanrpc

## TODO and Notes
* Better parsing handler
* Notify for unregistered methods, methods that are not registred, will be deleted from queue
* ...


## Example
```go
    // beanrpc
    package main

    import (
	    "github.com/MarinX/beanrpc"
	    "log"
	    "time"
    )

    type Test struct {
	    Name string
	    Age  int
    }

    func main() {

	    r := beanrpc.New("localhost:11300")

	    //opens tube for procesing
	    if err := r.Open("mytube"); err != nil {
		    log.Println(err)
		    return
	    }

	    //register method

	    r.On("mymethod", func(c *beanrpc.Context) {

		    log.Println("Buffered output->", string(c.Buff()))

		    log.Println("Job id->", c.Id())

		    //bind your type
		    var params string

		    if err := c.Bind(&params); err != nil {
			    log.Println(err)
		    }

		    log.Println("Params->", params)
	    })

	    r.On("secondmethod", func(c *beanrpc.Context) {
		    log.Println("Second method called!")

		    //can bind structs also
		    var str Test
		    c.Bind(&str)
		    log.Println(str)

	    })

	    go PushJobs(r)

	    //blocking method!
	    r.Run()

	    /*
		    OUTPUT:
			    2015/07/26 18:00:09 Buffered output->         {"Method":"mymethod","Params":"HelloWorld"}
			    2015/07/26 18:00:09 Job id-> 244
			    2015/07/26 18:00:09 Params-> HelloWorld
			    2015/07/26 18:00:09 Second method called!
			    2015/07/26 18:00:09 {Marin 22}
	    */
    }

    func PushJobs(r *beanrpc.BeanWorker) {
	    time.Sleep(2 * time.Second)

	    r.Put("mymethod", "HelloWorld", 1)

	    r.Put("secondmethod", Test{
		    Name: "Marin",
		    Age:  22,
	    }, 1)
    }
```

## License
This library is under the MIT License

## Author
Marin Basic 
