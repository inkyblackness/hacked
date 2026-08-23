#pragma once

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * Codepoint values must be Unicode scalar values: @code 0x0000..0xD7FF@endcode or @code 0xE000..0x10FFFF@endcode.
 * Surrogate code points and values above @code 0x10FFFF@endcode are invalid.
 */
typedef uint32_t Codepoint;

typedef enum
{
   TEXT_CHARSET_UTF8,
   TEXT_CHARSET_CODEPAGE437
} TextCharset;

typedef enum [[nodiscard]]
{
   /**
    * Processing could finish as part of the call.
    */
   TEXT_ENCODING_DONE,
   /**
    * Some of the provided input data can not be processed.
    */
   TEXT_ENCODING_INPUT_INVALID,
   /**
    * The provided output size is not enough to fully consume input.
    */
   TEXT_ENCODING_MORE_OUTPUT_NEEDED,
   /**
    * The provided input is incomplete and will require further data to create output.
    */
   TEXT_ENCODING_MORE_INPUT_NEEDED
} TextEncodingResult;

/**
 * An opaque decoder structure.
 */
typedef struct TextDecoder TextDecoder;

/**
 * Request to create a decoder instance.
 *
 * The returned decoder instance must be released with @ref textDecoderRelease "textDecoderRelease".
 *
 * @param charset the charset for which to create a decoder
 * @return the corresponding instance, or @code NULL@endcode if given charset is not supported
 */
[[nodiscard]] extern TextDecoder *textDecoderCreate(TextCharset charset);

/**
 * Releases the decoder instance. In case it was a dynamically allocated memory, it will be freed.
 *
 * @param decoder address of the decoder pointer to release; the pointed-to-pointer will be set to @code NULL@endcode.
 */
extern void textDecoderRelease(TextDecoder **decoder);

/**
 * Return the maximum amount of possible output codepoints per input byte.
 *
 * @param decoder the decoder to query
 * @return the maximum for the given decoder
 */
[[nodiscard]] extern size_t textDecoderMaxOutputPerInputByte(TextDecoder const *decoder);

/**
 * Resets the decoder instance to a state as if it was freshly initialized.
 * Any previous pending data in the internal state is dropped.
 *
 * @param decoder the decoder instance to reset
 */
extern void textDecoderReset(TextDecoder *decoder);

/**
 * Decode bytes and provide codepoints.
 *
 * The decoder will read as many input bytes as it can and produces output codepoints from them.
 *
 * The caller needs to interpret the output parameters according to the return value:
 * - TEXT_ENCODING_DONE: all provided input bytes could be consumed and output was generated.
 *   @code*outSize@endcode will specify how many codepoints were written to out.
 *   @code*inSize@endcode will be equal to its input value.
 *   The caller can continue with the next block.
 *
 * - TEXT_ENCODING_INPUT_INVALID: in contained at least one input byte that could not be interpreted.
 *   @code*outSize@endcode will specify how many codepoints were written to out until the problem occurred.
 *   @code*inSize@endcode will specify how many bytes were consumed and points to the first byte that can not be interpreted.
 *   The caller needs to perform some error correction. Strategies can include: abort decoding,
 *   skip one input byte and retry, or reset the decoder and retry. For selecting the proper strategy,
 *   the caller should take into account the previous return value. If the previous call returned
 *   TEXT_ENCODING_MORE_INPUT_NEEDED, then the decoder may have a pending state that may not be resolved,
 *   even when skipping bytes.
 *
 * - TEXT_ENCODING_MORE_INPUT_NEEDED: input data could be read, yet not everything could create codepoints.
 *   @code*outSize@endcode will specify how many codepoints were written to out.
 *   @code*inSize@endcode will be equal to its input value.
 *   The caller should continue with the next block that follows the pending data.
 *
 * - TEXT_ENCODING_MORE_OUTPUT_NEEDED: out data is exhausted, not all input could be consumed.
 *   @code*outSize@endcode will specify how many codepoints were written to out and will be equal to its input value.
 *   @code*inSize@endcode will specify how many bytes were consumed and points to the first byte the caller should provide again.
 *   The caller should continue with a new call and a fresh output buffer, having input continue where the last call left off.
 *
 * @param decoder the decoder instance
 * @param out pointer to output codepoints to write to. Must not be @code NULL@endcode.
 * @param outSize in/out parameter. Input: number of available out items to write to, at least 1. Output: See return value
 * @param in pointer to input bytes to decode. Must not be @code NULL@endcode.
 * @param inSize in/out parameter. Input: number of provided bytes, at least 1. Output: See return value
 * @return the result of the processing
 */
[[nodiscard]] extern TextEncodingResult textDecoderDecode(TextDecoder *decoder, Codepoint *out, size_t *outSize, uint8_t const *in, size_t *inSize);

/**
 * Flushes the decoder. It will consume any pending input data from its state and provide any resulting codepoints.
 * Once successfully flushed, the decoder is in initial state as if it were freshly constructed.
 *
 * The caller needs to interpret the output parameter according to the return value:
 * - TEXT_ENCODING_DONE: all pending input bytes could be consumed and output was maybe generated.
 *   @code*outSize@endcode will specify how many codepoints were written to out, which may be zero.
 *   The caller can use the decoder for a new text.
 *
 * - TEXT_ENCODING_MORE_INPUT_NEEDED: pending input data is incomplete and output may or may not have been generated.
 *   @code*outSize@endcode will specify how many codepoints were written to out, which may be zero.
 *   The caller needs to perform some error correction. Strategies can include: abort decoding and reset the decoder,
 *   or provide the missing input data.
 *
 * - TEXT_ENCODING_MORE_OUTPUT_NEEDED: out data is exhausted, not all input could be consumed.
 *   @code*outSize@endcode will specify how many codepoints were written to out and will be equal to its input value.
 *   The caller should continue with a new call and a fresh output buffer.
 *
 * @param decoder the decoder instance
 * @param out pointer to output codepoints to write to. Must not be @code NULL@endcode.
 * @param outSize in/out parameter. Input: number of available out items to write to, at least 1. Output: See return value
 * @return the result of the flushing
 */
[[nodiscard]] extern TextEncodingResult textDecoderFlush(TextDecoder *decoder, Codepoint *out, size_t *outSize);

/**
 * An opaque encoder structure.
 */
typedef struct TextEncoder TextEncoder;

/**
 * Request to create an encoder instance.
 *
 * The returned encoder instance must be released with @ref textEncoderRelease "textEncoderRelease".
 *
 * The returned encoder will issue '?' in case the input codepoint can not be represented in the
 * target encoding.
 *
 * @param charset the charset for which to create an encoder
 * @return the corresponding instance, or @code NULL@endcode if given charset is not supported
 */
[[nodiscard]] extern TextEncoder *textEncoderCreate(TextCharset charset);

/**
 * Releases the encoder instance. In case it was a dynamically allocated memory, it will be freed.
 *
 * @param encoder address of the encoder pointer to release; the pointed-to-pointer will be set to @code NULL@endcode.
 */
extern void textEncoderRelease(TextEncoder **encoder);

/**
 * Return the maximum amount of possible output bytes per input codepoint.
 *
 * @param encoder the encoder to query
 * @return the maximum for the given encoder
 */
[[nodiscard]] extern size_t textEncoderMaxOutputPerCodepoint(TextEncoder const *encoder);

/**
 * Resets the encoder instance to a state as if it was freshly initialized.
 * Any previous pending data in the internal state is dropped.
 *
 * @param encoder the encoder to reset
 */
extern void textEncoderReset(TextEncoder *encoder);

/**
 * Encode codepoints to a byte buffer.
 *
 * It is encoder-specific how to handle a codepoint that can not be represented by the encoding.
 * It could skip it, use a placeholder character, or treat it as invalid input; all depending on how
 * the encoder was created.
 *
 * The caller needs to interpret the output parameters according to the return value:
 * - TEXT_ENCODING_DONE: all provided input codepoints could be consumed and output was generated.
 *   @code*outSize@endcode will specify how many bytes were written to out.
 *   @code*inSize@endcode will be equal to its input value.
 *   The caller can continue with the next block.
 *
 * - TEXT_ENCODING_INPUT_INVALID: in contained at least one codepoint that could not be interpreted.
 *   @code*outSize@endcode will specify how many bytes were written to out until the problem occurred.
 *   @code*inSize@endcode will specify how many codepoints were consumed and points to the first codepoint that can not be interpreted.
 *   The caller needs to perform some error correction. Strategies can include: abort encoding,
 *   skip one input codepoint and retry, or reset the encoder and retry. For selecting the proper strategy,
 *   the caller should take into account the previous return value. If the previous call returned
 *   TEXT_ENCODING_MORE_INPUT_NEEDED, then the encoder may have a pending state that may not be resolved,
 *   even when skipping codepoints.
 *
 * - TEXT_ENCODING_MORE_INPUT_NEEDED: input data could be read, yet not everything could create bytes.
 *   @code*outSize@endcode will specify how many bytes were written to out.
 *   @code*inSize@endcode will be equal to its input value.
 *   The caller should continue with the next block that follows the pending data.
 *   Note that for encoding, this is a seldom case and only for special encoders.
 *
 * - TEXT_ENCODING_MORE_OUTPUT_NEEDED: out data is exhausted, not all input could be consumed.
 *   @code*outSize@endcode will specify how many bytes were written to out and will be equal to its input value.
 *   @code*inSize@endcode will specify how many codepoints were consumed and points to the first codepoint the caller should provide again.
 *   The caller should continue with a new call and a fresh output buffer, having input continue where the last call left off.
 *
 * @param encoder the encoder instance
 * @param out pointer to output bytes to write to. Must not be @code NULL@endcode.
 * @param outSize in/out parameter. Input: number of available bytes to write to, at least 1. Output: See return value.
 * @param in pointer to input codepoints to encode. Must not be @code NULL@endcode.
 * @param inSize in/out parameter. Input: number of available codepoints to encode, at least 1. Output: See return value.
 * @return the result of the processing.
 */
[[nodiscard]] extern TextEncodingResult textEncoderEncode(TextEncoder *encoder, uint8_t *out, size_t *outSize, Codepoint const *in, size_t *inSize);

/**
 * Flushes the encoder. It will consume any pending input data from its state and provide any resulting bytes.
 * Once successfully flushed, the encoder is in initial state as if it were freshly constructed.
 *
 * The caller needs to interpret the output parameter according to the return value:
 * - TEXT_ENCODING_DONE: all pending input codepoints could be consumed and output was maybe generated.
 *   @code*outSize@endcode will specify how many bytes were written to out, which may be zero.
 *   The caller can use the encoder for a new text.
 *
 * - TEXT_ENCODING_MORE_INPUT_NEEDED: pending input codepoints are incomplete and output may or may not have been generated.
 *   @code*outSize@endcode will specify how many bytes were written to out, which may be zero.
 *   The caller needs to perform some error correction. Strategies can include: abort encoding and reset the encoder,
 *   or provide the missing input data.
 *
 * - TEXT_ENCODING_MORE_OUTPUT_NEEDED: out data is exhausted, not all input could be consumed.
 *   @code*outSize@endcode will specify how many bytes were written to out and will be equal to its input value.
 *   The caller should continue with a new call and a fresh output buffer.
 *
 * @param encoder the encoder instance
 * @param out pointer to output bytes to write to
 * @param outSize in/out parameter. Input: number of available out items to write to. Output: See return value
 * @return the result of the flushing. See documentation of TextEncodingResult for details.
 */
[[nodiscard]] extern TextEncodingResult textEncoderFlush(TextEncoder *encoder, uint8_t *out, size_t *outSize);

#ifdef __cplusplus
}
#endif
