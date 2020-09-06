package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

const escapeChar = '\\'
const zeroByte rune = 0

var ErrInvalidString = errors.New("invalid string")

func Unpack(packed string) (string, error) {
	// fmt.Printf("-> '%s'\n", packed)
	res := &strings.Builder{}
	repeatCount := 1
	lastIndex := len(packed) - 1
	var prevChar rune

	var escapeMode, repeatMode bool

	add := func(char rune) {
		addChar(res, prevChar, repeatCount)
		prevChar = char
		repeatCount = 1
		repeatMode, escapeMode = false, false
	}

	for i, char := range packed {
		// fmt.Printf("%d %c : prev=%c : %s\n", i, char, prevChar, res.String())

		switch {
		case char == escapeChar:
			if escapeMode {
				add(char)
				continue
			}
			if i == lastIndex {
				// незаэскейпленный концевой эскейп: abcd\
				return "", ErrInvalidString
			}
			escapeMode = true

		case unicode.IsDigit(char):
			if escapeMode {
				add(char)
				continue
			}
			if i == 0 || repeatMode {
				// цифра в начале строки или цифра после незаэскейпленной цифры 5abcd, abc35q
				return "", ErrInvalidString
			}
			repeatCount = int(char - '0')
			repeatMode = true

		default:
			// не цифра не слеш
			if escapeMode {
				// некорректный эскейп: abc\xde или \x, при этом abc\\xde - корректен
				return "", ErrInvalidString
			}
			add(char)
		}
		// fmt.Printf("%d %c : prev=%c : flags=%d : %s\n", i, char, prevChar, flags, res.String())
	}
	// поскольку в цикле всегда добавляли предыдущий символ (с некоторыми оговорками),
	// то последний prevChar не был добавлен
	addChar(res, prevChar, repeatCount)
	// fmt.Printf("%s -> %s\n", packed, res.String())
	return res.String(), nil
}

func addChar(buff *strings.Builder, char rune, repeat int) {
	if char != zeroByte {
		for i := 0; i < repeat; i++ {
			buff.WriteRune(char)
		}
	}
}
