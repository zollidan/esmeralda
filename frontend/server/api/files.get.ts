export default defineEventHandler(async (event) => {
  const runtimeConfig = useRuntimeConfig();

  try {
    const response = await fetch(`${runtimeConfig.apiBaseURL}/api/v1/files`);

    if (!response.ok) {
      return { error: "Failed to fetch files" };
    }

    const data = await response.json();

    return data;
  } catch (error) {
    return { error: "An error occurred while fetching parsers" };
  }
});
