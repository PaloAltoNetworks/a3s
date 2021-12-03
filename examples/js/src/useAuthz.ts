interface UseAuthzOptions {
  /**
   * The base URL for the a3s backend. Shouldn't include the `/` in the end.
   * Example: `https://127.0.0.1:44443`
   */
  baseUrl: string;
}

export const useAuthz =
  ({ baseUrl }: UseAuthzOptions) =>
  ({ resource, token, namespace }: { resource: string; token: string; namespace: string }) =>
    fetch(`${baseUrl}/authz`, {
      method: "POST",
      body: JSON.stringify({
        token,
        action: "GET",
        resource,
        namespace,
        audience: baseUrl,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    }).then(
      (res) => res.status === 204,
      (e) => {
        window.alert(`Unknown error ${e}`);
        return false;
      }
    );
