// Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/asm/asm.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package asm provides support for dealing with EVM assembly instructions (e.g., disassembling them).

# Source Files

Each file provides the following features.
  - asm.go: provides instruction iterators for EVM assembly instructions.
  - compiler.go: provides a compiler which compiles input tokens and returns EVM binaries.
  - lexer.go: provides the basic construct for parsing source code and turning them into tokens.
*/
package asm
