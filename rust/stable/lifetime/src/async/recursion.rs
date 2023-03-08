use futures::future::LocalBoxFuture;

struct NonCopyClone<'bytes> {
    bytes: &'bytes [u8],
}

fn recursion<'bytes, 'data>(
    data: &'data NonCopyClone<'bytes>,
) -> LocalBoxFuture<'bytes, &'data NonCopyClone<'bytes>>
where
    'data: 'bytes,
{
    Box::pin(async move {
        if data.bytes.is_empty() {
            return data;
        }

        data
    })
}
