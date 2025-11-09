export default defineEventHandler(async (event) => {
  const runtimeConfig = useRuntimeConfig();

  const body = await readBody(event);
  const { name, description } = body;

  try {
    const response = await fetch(`${runtimeConfig.apiBaseURL}/api/v1/parsers`, {
      method: "POST",
      body: JSON.stringify({ name, description }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      return { error: "Failed to create parser" };
    }

    const data = await response.json();

    return data;
  } catch (error) {
    return { error: "An error occurred while fetching parsers" };
  }
});
