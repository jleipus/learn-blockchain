package merkel_test

import (
	"crypto/sha256"
	"testing"

	"github.com/jleipus/learn-blockchain/internal/merkel"
	"github.com/stretchr/testify/assert"
)

func TestNewMerkelTree(t *testing.T) {
	t.Run("single element", func(t *testing.T) {
		data := [][]byte{[]byte("single")}
		expectedHash := sha256.Sum256(data[0])
		expectedRootHash := sha256.Sum256(append(expectedHash[:], expectedHash[:]...))

		tree := merkel.NewTree(data)

		assert.NotNil(t, tree)
		assert.NotNil(t, tree.Root)
		assert.NotNil(t, tree.Root.Left)
		assert.NotNil(t, tree.Root.Right)
		assert.Equal(t, expectedHash[:], tree.Root.Left.GetData())
		assert.Equal(t, expectedHash[:], tree.Root.Right.GetData())
		assert.Equal(t, expectedRootHash[:], tree.Root.GetData())
	})

	t.Run("two elements", func(t *testing.T) {
		data := [][]byte{
			[]byte("first"),
			[]byte("second"),
		}

		expectedLeftHash := sha256.Sum256(data[0])
		expectedRightHash := sha256.Sum256(data[1])
		expectedRootHash := sha256.Sum256(append(expectedLeftHash[:], expectedRightHash[:]...))

		tree := merkel.NewTree(data)

		assert.NotNil(t, tree)
		assert.NotNil(t, tree.Root)
		assert.NotNil(t, tree.Root.Left)
		assert.NotNil(t, tree.Root.Right)
		assert.Equal(t, expectedLeftHash[:], tree.Root.Left.GetData())
		assert.Equal(t, expectedRightHash[:], tree.Root.Right.GetData())
		assert.Equal(t, expectedRootHash[:], tree.Root.GetData())
	})
}
