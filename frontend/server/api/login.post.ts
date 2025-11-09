export default defineEventHandler(async (event) => {
  const runtimeConfig = useRuntimeConfig();

  const body = await readBody(event);
  const { username, password } = body;

  try {
    const response = await fetch(
      `${runtimeConfig.apiBaseURL}/api/v1/auth/login`,
      {
        method: "POST",
        body: JSON.stringify({ username, password }),
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      return { error: "Failed to fetch parsers" };
    }

    const data = await response.json();

    return data;
  } catch (error) {
    return { error: "An error occurred while fetching parsers" };
  }
});
