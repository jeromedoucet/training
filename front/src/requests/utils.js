
export async function handleResponse (res) {
  const body = await _parseJSON(res)
  if (!res.ok) {
    return Promise.reject(body)
  } else {
    return body
  }
}

function _parseJSON (response) {
  return response.text().then(function (text) {
    return text ? JSON.parse(text) : {}
  })
}
