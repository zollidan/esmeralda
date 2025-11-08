export default defineEventHandler(async (event) => {
  const runtimeConfig = useRuntimeConfig();

  try {
    const response = await fetch(`${runtimeConfig.apiBaseURL}/api/v1/parsers`);

    if (!response.ok) {
      return { error: "Failed to fetch parsers" };
    }

    const data = await response.json();
  } catch (error) {
    return { error: "An error occurred while fetching parsers" };
  }
});
