#include <tox/tox.h>

/* Convenient macro:
 * Creates the C function to directly register a given callback */
#define CREATE_HOOK(x) \
static void set_##x(Tox *tox, void *t) { \
  tox_##x(tox, hook_##x, t); \
}

// Headers for the exported GO functions in hooks.go
void hook_callback_self_connection_status(Tox*, TOX_CONNECTION, void*);
void hook_callback_friend_name(Tox*, uint32_t, const uint8_t*, size_t, void*);
void hook_callback_friend_status_message(Tox*, uint32_t, const uint8_t*, size_t, void*);
void hook_callback_friend_status(Tox*, uint32_t, TOX_USER_STATUS, void*);
void hook_callback_friend_connection_status(Tox*, uint32_t, TOX_CONNECTION, void*);
void hook_callback_friend_typing(Tox*, uint32_t, bool, void*);
void hook_callback_friend_read_receipt(Tox*, uint32_t, uint32_t, void*);
void hook_callback_friend_request(Tox*, const uint8_t*, const uint8_t*, size_t, void*);
void hook_callback_friend_message(Tox*, uint32_t, TOX_MESSAGE_TYPE, const uint8_t*, size_t, void*);
void hook_callback_file_recv_control(Tox*, uint32_t, uint32_t, TOX_FILE_CONTROL, void*);
void hook_callback_file_chunk_request(Tox*, uint32_t, uint32_t, uint64_t, size_t, void*);
void hook_callback_file_recv(Tox*, uint32_t, uint32_t, uint32_t, uint64_t, const uint8_t*, size_t, void*);
void hook_callback_file_recv_chunk(Tox*, uint32_t, uint32_t, uint64_t, const uint8_t*, size_t, void*);
void hook_callback_friend_lossy_packet(Tox*, uint32_t, const uint8_t*, size_t, void*);
void hook_callback_friend_lossless_packet(Tox*, uint32_t, const uint8_t*, size_t, void*);

CREATE_HOOK(callback_self_connection_status)
CREATE_HOOK(callback_friend_name)
CREATE_HOOK(callback_friend_status_message)
CREATE_HOOK(callback_friend_status)
CREATE_HOOK(callback_friend_connection_status)
CREATE_HOOK(callback_friend_typing)
CREATE_HOOK(callback_friend_read_receipt)
CREATE_HOOK(callback_friend_request)
CREATE_HOOK(callback_friend_message)
CREATE_HOOK(callback_file_recv_control)
CREATE_HOOK(callback_file_chunk_request)
CREATE_HOOK(callback_file_recv)
CREATE_HOOK(callback_file_recv_chunk)
CREATE_HOOK(callback_friend_lossy_packet)
CREATE_HOOK(callback_friend_lossless_packet)
