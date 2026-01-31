package model

import (
	cryptorand "crypto/rand"
	"math/rand"
	"sync"
	"testing"
)

const (
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphabetLength = byte(len(alphabet))
)

type StorageTest struct {
	length int
}

// Оригинальная версия
func (s *StorageTest) randomStringOld(letters string) string {
	b := make([]byte, s.length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Оптимизированная версия
func (s *StorageTest) randomStringNew() string {
	b := make([]byte, s.length)
	if _, err := cryptorand.Read(b); err != nil {
		return ""
	}
	for i := range b {
		b[i] = alphabet[b[i]%alphabetLength]
	}
	return string(b)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 8)
		return &buf
	},
}

// версия с pool
func (s *StorageTest) randomStringPooled() string {
	bufPtr := bufPool.Get().(*[]byte)
	buf := *bufPtr
	if len(buf) < s.length {
		buf = make([]byte, s.length)
		*bufPtr = buf
	}
	defer bufPool.Put(bufPtr)

	_, _ = cryptorand.Read(buf[:s.length])
	for i := range s.length {
		buf[i] = alphabet[buf[i]%alphabetLength]
	}
	return string(buf[:s.length])
}

// Бенчмарк: старая версия
func BenchmarkRandomString_Old(b *testing.B) {
	storage := &StorageTest{length: 8}
	letters := alphabet
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = storage.randomStringOld(letters)
	}
}

// Бенчмарк: новая версия
func BenchmarkRandomString_New(b *testing.B) {
	storage := &StorageTest{length: 8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = storage.randomStringNew()
	}
}

func BenchmarkRandomString_Pooled(b *testing.B) {
	storage := &StorageTest{length: 8}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = storage.randomStringPooled()
	}
}
