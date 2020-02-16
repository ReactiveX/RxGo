# RxGo
[![Join the chat at https://gitter.im/ReactiveX/RxGo](https://badges.gitter.im/ReactiveX/RxGo.svg)](https://gitter.im/ReactiveX/RxGo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/ReactiveX/RxGo.svg?branch=master)](https://travis-ci.org/ReactiveX/RxGo)
[![Coverage Status](https://coveralls.io/repos/github/ReactiveX/RxGo/badge.svg?branch=master)](https://coveralls.io/github/ReactiveX/RxGo?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/reactivex/rxgo)](https://goreportcard.com/report/github.com/reactivex/rxgo)

Reactive Extensions for the Go Language

## Getting Started

[ReactiveX](http://reactivex.io/), or Rx for short, is an API for programming with observable streams. This is the official ReactiveX API for the Go language.

**ReactiveX** is a new, alternative way of asynchronous programming to callbacks, promises and deferred. It is about processing streams of events or items, with events being any occurrences or changes within the system. A stream of events is called an [observable](http://reactivex.io/documentation/contract.html).

An operator is basically a function that defines an observable, how and when it should emit data. The list of operators covered is available [here](README.md#operators).

## RxGo

The **RxGo** implementation is based on the idea of [pipelines](https://blog.golang.org/pipelines). In a nutshell, a pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function.

Let's check at a concrete example with each blue box being an operator:
* We create a static observable based on a fixed list of items using the `Just` operator.
* We define a transformation function using the `Map` operator (a circle into a square).
* We filter each red square using the `Filter` operator.

![](res/rx.png)

This stream produced in the target channel two items (a yellow and a green square).
Each operator is a transformation stage connected by channels. By default, everything is sequential. Yet, we can easily leverage modern CPU architectures by defining multiple instances of the same operator (each operator instance being a goroutine connected to a common channel).

## Operators

## Call for Maintainers
The development of RxGo v2 has started (`v2` branch). Ongoing discussions can be found in [#99](https://github.com/ReactiveX/RxGo/issues/99). 

We are welcoming anyone who's willing to help us maintaining or would like to become a core developer to get in touch with us.

## Contributions
All contributions are welcome, both in development and documentation! Be sure you check out [contributions](https://github.com/ReactiveX/RxGo/wiki/Contributions) and [roadmap](https://github.com/ReactiveX/RxGo/wiki/Roadmap).
