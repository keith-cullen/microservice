async function callV1Get(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }
    const data = await response.json();
    console.log(data);
    document.body.innerHTML =
      "<p>Request: " + `${url}` + "</p>" +
      "<p>Response: " + `${data}` + "</p>"
  } catch (error) {
    console.error(error.message);
  }
}

const url = "https://localhost/v1/get?name=Bob"
callV1Get(url);
