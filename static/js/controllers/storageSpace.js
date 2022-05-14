"use strict";

/*
example 200 OK response:
{
  "all":931.51,
  "free":36.4,
  "used":895.11,
  "usedPercentage":96.09
}
*/
export async function getStorageSpaceOrNull() {
  // this API call is expected to fail on unsupported OSes, so we want to gracefully
  // fail on errors.
  return fetch(`/api/storage-space`, {
    method: "GET",
    credentials: "include",
  })
  .then((response) => {
    if (!response.ok) {
      return Promise.resolve(null);
    }
    return Promise.resolve(response.json());
  })
  .catch((error) => {
    return Promise.resolve(null);
  });
}
