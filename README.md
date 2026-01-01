# go-jobs 🚀  
**A production-ready, cancellable background job system for Go**

[![Go Reference](https://pkg.go.dev/badge/github.com/chokey2nv/go-jobs.svg)](https://pkg.go.dev/github.com/chokey2nv/go-jobs)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

---

## Overview

`go-jobs` is a lightweight, storage-agnostic job execution library for Go.
It is designed to solve **long-running backend tasks** (AI calls, document processing,
report generation, imports, exports) without HTTP / GraphQL timeouts.

The library emphasizes:
- Explicit control
- Predictable execution
- Clean cancellation
- No external infrastructure required

---

## Why go-jobs?

Running long tasks directly inside API handlers leads to:
- Request timeouts
- Poor user experience
- Difficult retries and recovery

`go-jobs` solves this by executing work asynchronously while exposing
a **queryable job state** that clients can poll.

---

## Features

- ✅ Asynchronous job execution
- ✅ Worker pool with backpressure
- ✅ Job progress reporting
- ✅ Job cancellation via context
- ✅ Storage injection (memory, DB, Redis, etc.)
- ✅ Safe concurrency
- ✅ GraphQL / REST friendly
- ✅ Test-first design

---

## Job Lifecycle


Each job tracks:
- Status
- Progress (0–100)
- Message
- Result
- Error
- CreatedAt / UpdatedAt

---

## Installation

```bash
go get github.com/chokey2nv/go-jobs

store := stores.NewMemoryStore()
svc   := job.New(store)


j, err := svc.StartAsync(ctx, "generate-report",
	func(ctx context.Context, report types.ProgressReporter) (any, error) {

		report.Progress(10, "Starting")
		time.Sleep(time.Second)

		report.Progress(50, "Processing")
		time.Sleep(time.Second)

		report.Progress(100, "Done")
		return "success", nil
	},
)



job, _ := svc.Get(ctx, j.ID)

fmt.Println(job.Status)
fmt.Println(job.Progress)
fmt.Println(job.Result)

// cancel task 

err := svc.Cancel(ctx, j.ID)


```
