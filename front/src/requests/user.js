import { handleResponse } from '@/requests/utils'

export function createUser (form) {
  return fetch('/app/public/users', {
    method: 'POST',
    credentials: 'same-origin',
    body: JSON.stringify(form)
  })
    .then(handleResponse)
}

export function login (form) {
  return fetch('/app/public/login', {
    method: 'POST',
    credentials: 'same-origin',
    body: JSON.stringify(form)
  })
    .then(handleResponse)
}

export function logout () {
  return fetch('/app/public/logout', {
    method: 'POST',
    credentials: 'same-origin'
  })
    .then(handleResponse)
}
