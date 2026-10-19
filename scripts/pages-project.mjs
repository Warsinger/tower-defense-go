export async function pagesProject(fetcher, base, headers, allowCreate) {
  const project = 'tower-def';
  let operation = 'lookup';
  let response = await fetcher(`${base}/${project}`, { headers });
  if (response.status === 404 && allowCreate) {
    operation = 'create';
    response = await fetcher(base, {
      method: 'POST', headers: { ...headers, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: project, production_branch: 'main' }),
    });
  }
  if (!response.ok) {
    const detail = await response.json().catch(() => ({}));
    const codes = (detail.errors || []).map(error => error.code).filter(code => Number.isInteger(code)).join(",");
    throw new Error(`Cloudflare project ${operation} failed: HTTP ${response.status}; error codes: ${codes || "unavailable"}`);
  }
  const result = await response.json();
  if (!result.success || result.result?.name !== project || result.result.production_branch !== 'main' || result.result.source) {
    throw new Error('Expected the tower-def Direct Upload project with production branch main');
  }
  return result.result;
}
