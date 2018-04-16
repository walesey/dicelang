# Dicelang

[![Build Status](https://travis-ci.org/walesey/dicelang.svg?branch=master)](https://travis-ci.org/walesey/dicelang)

### A programming language for calculating dice roll statistics.

## Installation

` $ go get github.com/walesey/dicelang `

## Usage

### Roll some dice
`$ dicelang "resolve 2d6.add" `

### calculate the mean value for rolling 2 dice and adding the result
`$ dicelang "mean 2d6.add" `

### generate a histogram
`$ dicelang "hist 2d6.add" `

### adding multiple different dice
`$ dicelang "hist [2d6.add, d10.add, d6.add]" `

### rolling a number of dice equal to a previous roll
`$ dicelang "hist 2d6.add.d6.add" `

### number of dice that roll equal to or greater than a value
`$ dicelang "hist 10d6.5+" `

![Example](example.png?raw=true)

