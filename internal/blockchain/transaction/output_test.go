package transaction_test

import (
	"testing"

	"github.com/jleipus/learn-blockchain/internal/blockchain/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerializeDeserialize(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		o := &transaction.TxOutput{}
		serialized := o.Serialize()
		require.NotEmpty(t, serialized)

		var deserialized transaction.TxOutput
		err := deserialized.Deserialize(serialized)
		require.NoError(t, err)

		assert.Equal(t, o, &deserialized)
	})

	t.Run("with basic data", func(t *testing.T) {
		o := &transaction.TxOutput{
			Value:      100,
			PubKeyHash: []byte("test-pubkey"),
		}
		serialized := o.Serialize()
		require.NotEmpty(t, serialized)

		var deserialized transaction.TxOutput
		err := deserialized.Deserialize(serialized)
		require.NoError(t, err)

		assert.Equal(t, o, &deserialized)
	})

	t.Run("multiple outputs", func(t *testing.T) {
		outputs := []transaction.TxOutput{
			{Value: 100, PubKeyHash: []byte("pubkey1")},
			{Value: 200, PubKeyHash: []byte("pubkey2")},
		}

		serialized, err := transaction.SerializeOutputs(outputs)
		require.NoError(t, err)
		require.NotEmpty(t, serialized)

		deserializedOutputs, err := transaction.DeserializeOutputs(serialized)
		require.NoError(t, err)

		require.Len(t, deserializedOutputs, 2)
		assert.Equal(t, outputs[0], deserializedOutputs[0])
		assert.Equal(t, outputs[1], deserializedOutputs[1])
	})
}
