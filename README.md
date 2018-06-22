# Mental Poker Algorithm Evaluation

This project houses an evaluation of the [Mental Poker](https://en.wikipedia.org/wiki/Mental_poker) algorithm using SRA
commutative encryption written in Go. There are two packages here that are written to be simple to read: [sra](sra)
which is the SRA implementation that generates keys and does encryption/decryption, and [deck](deck) which is the
shuffling and card drawing algorithm.

This is meant to be a demo more than a library. Feel free to copy/reuse any code or learn from it (MIT licensed). With
`GOPATH` set, the code and all dependencies can be fetched with `go get -u github.com/cretz/go-mental-poker/...`.

## Overview

The SRA keys are two numbers based on a shared prime that, when one is applied to a value, the other can be used to
"unapply" (i.e. encrypt and decrypt respectively). It is commutative, meaning that multiple encryptions can occur on a
value in any order to generate an encrypted value. Then the counteracting decryptions can occur on the encrypted value,
in any order, and the original value will be the result.

The mental poker algorithm provides a way for disparate players to shuffle a set of cards without a trusted third party.
Then it allows decrypting cards only by the player drawing. Here's how shuffling works:

* A set of cards is decided to be the deck and a prime number is generated and shared
* Starting with the unencrypted deck, the set of cards is sent to each player serially where that player encrypts each
  card with a single SRA key they create, then shuffles the cards before passing the card set to the next player
* Each card in the deck has now been encrypted (on top of each other) by every player and re-shuffled by every player
* Now the deck is sent back around, and each player decrypts every card with the key they encrypted before, and
  re-encrypts them with new keys this time making sure the keys are different for each card
* Finally, the completed deck is sent back around so each player can map their per-card keys to the final encrypted
  values

Here's how drawing works:

* Some event, recognized as legitimate by all players, occurs where a player draws
* That player asks all the other players for their decryption key for that card
* The player then uses their decryption keys + the player's own for that card to get the actual card value

At the end of the game, all cards and decryption keys should be made visible so each player can verify that all cards
were handled properly.

An example of this can be seen in [deck/deck_test.go](deck/deck_test.go) where a regular 52-card deck is shuffled and 3
players are given 7 cards. When at the root of the repository, run the following:

    go test -v ./deck

The output will look something like:

    === RUN   TestSimpleDraw
    All cards:          [2♠ 2♥ 2♦ 2♣ 3♠ 3♥ 3♦ 3♣ 4♠ 4♥ 4♦ 4♣ 5♠ 5♥ 5♦ 5♣ 6♠ 6♥ 6♦ 6♣ 7♠ 7♥ 7♦ 7♣ 8♠ 8♥ 8♦ 8♣ 9♠ 9♥ 9♦ 9♣ 10♠ 10♥ 10♦ 10♣ J♠ J♥ J♦ J♣ Q♠ Q♥ Q♦ Q♣ K♠ K♥ K♦ K♣ A♠ A♥ A♦ A♣]
    Deck after shuffle: [10♠ 9♥ A♣ 10♥ A♠ 4♥ K♣ 5♥ 10♦ 7♠ Q♠ K♥ 3♣ 7♣ 9♣ 9♦ K♦ 6♥ 3♥ J♥ 6♣ 2♠ 8♠ J♠ A♦ 4♦ J♣ 8♥ Q♦ 9♠ 5♦ 2♣ 2♥ 3♠ Q♥ K♠ J♦ 7♥ 8♣ 10♣ A♥ 6♠ 5♣ 2♦ 4♣ Q♣ 5♠ 6♦ 3♦ 8♦ 4♠ 7♦]
    Deck after draws:   [10♠ 9♥ A♣ 10♥ A♠ 4♥ K♣ 5♥ 10♦ 7♠ Q♠ K♥ 3♣ 7♣ 9♣ 9♦ K♦ 6♥ 3♥ J♥ 6♣ 2♠ 8♠ J♠ A♦ 4♦ J♣ 8♥ Q♦ 9♠ 5♦]
    Alice's draw:       [7♦ 3♦ Q♣ 5♣ 10♣ J♦ 3♠]
    Bob's draw:         [4♠ 6♦ 4♣ 6♠ 8♣ K♠ 2♥]
    Ted's draw:         [8♦ 5♠ 2♦ A♥ 7♥ Q♥ 2♣]
    --- PASS: TestSimpleDraw (0.05s)
    PASS
    ok      github.com/cretz/go-mental-poker/deck   0.156s

The code in that test is the clearest overview of what the algorithm does.

## Benchmarks

One of the problems with this approach is speed. There are benchmarks for shuffling at
[deck/deck_bench_test.go](deck/deck_bench_test.go) and SRA key generation at
[sra/sra_bench_test.go](sra/sra_bench_test.go). In addition to the test mentioned previously, the benchmarks can be run
with the following at the repo root:

    go test ./... -bench=.

Excluding the output of the aforementioned test, here is the output on my mediocre Windows laptop:

    goos: windows                                                                   
    goarch: amd64                                                                   
    pkg: github.com/cretz/go-mental-poker/deck                                      
    BenchmarkShuffle2Players32Bits52Cards-8               10         154511300 ns/op
    BenchmarkShuffle2Players64Bits52Cards-8               10         183588590 ns/op
    BenchmarkShuffle3Players32Bits52Cards-8                5         227003880 ns/op
    BenchmarkShuffle3Players64Bits52Cards-8                5         286361780 ns/op
    BenchmarkShuffle6Players32Bits52Cards-8                3         456213900 ns/op
    BenchmarkShuffle6Players64Bits52Cards-8                2         588945400 ns/op
    BenchmarkShuffle2Players32Bits104Cards-8               5         310618440 ns/op
    BenchmarkShuffle2Players64Bits104Cards-8               3         381682066 ns/op
    BenchmarkShuffle3Players32Bits104Cards-8               3         444515500 ns/op
    BenchmarkShuffle3Players64Bits104Cards-8               2         561783900 ns/op
    BenchmarkShuffle6Players32Bits104Cards-8               1        1007646100 ns/op
    BenchmarkShuffle6Players64Bits104Cards-8               1        1440821900 ns/op
    PASS                                                                            
    ok      github.com/cretz/go-mental-poker/deck   27.004s                         
    goos: windows                                                                   
    goarch: amd64                                                                   
    pkg: github.com/cretz/go-mental-poker/sra                                       
    BenchmarkGenerateKeyPairSmallPrime32Bit-8          10000            121169 ns/op
    BenchmarkGenerateKeyPairSmallPrime64Bit-8           5000            317644 ns/op
    BenchmarkGenerateKeyPairSmallPrime128Bit-8          2000            599163 ns/op
    BenchmarkGenerateKeyPairMediumPrime32Bit-8         10000            104124 ns/op
    BenchmarkGenerateKeyPairMediumPrime64Bit-8          5000            325079 ns/op
    BenchmarkGenerateKeyPairMediumPrime128Bit-8         2000            770574 ns/op
    BenchmarkGenerateKeyPairLargePrime32Bit-8          10000            128342 ns/op
    BenchmarkGenerateKeyPairLargePrime64Bit-8           5000            317865 ns/op
    BenchmarkGenerateKeyPairLargePrime128Bit-8          2000            786605 ns/op
    BenchmarkGenerateKeyPairLargePrime256Bit-8           500           2831531 ns/op
    BenchmarkGenerateKeyPairLargePrime512Bit-8           100          18149114 ns/op
    BenchmarkGenerateKeyPairLargePrime1024Bit-8           10         109888740 ns/op
    BenchmarkGenerateKeyPairLargePrime2048Bit-8            1        2040425400 ns/op
    PASS                                                                            
    ok      github.com/cretz/go-mental-poker/sra    21.619s                         

This shows a few things (that could have been inferred just by understanding the algorithm):

* The size of the shared prime doesn't make that much of a difference in speed
* Increasing the number of bits for the key increases the key generation time exponentially
* Increasing the number of players increases shuffling time linearly
* Increasing the number of cards increases shuffling linearly

I am no expert here so I could be misreading, or more likely that I need more parameters and a larger sample size (that
I could increase w/ something like `-benchtime=5s` instead of the `1s` default) to draw any reasonable conclusions.

But even with this simple benchmark, it takes over a second on my laptop to shuffle a deck of 104 cards for 6 players.
This is because besides the first key each player generates, each player has to generate a key for every card. And this
is all done locally, so once network overhead is added, this could be a several second process.

## WARNING

This is just evaluation code and this mental poker algorithm and SRA encryption is known to have some weaknesses
including leaking card info. I have not researched mitigations to these weaknesses when developing this proof of
concept.