export async function pagesProject(fetcher, base, headers, allowCreate) {
  const project = 'tower-defense-go';
  let response = await fetcher(`${base}/${project}`, { headers });
  if (response.status === 404 && allowCreate) {
    response = await fetcher(base, {
      method: 'POST', headers: { ...headers, 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: project, production_branch: 'main' }),
    });
  }
  if (!response.ok) throw new Error(`Cloudflare project access/setup failed: HTTP ${response.status}`);
  const result = await response.json();
  if (!result.success || result.result?.name !== project || result.result.production_branch !== 'main' || result.result.source) {
    throw new Error('Expected the tower-defense-go Direct Upload project with production branch main');
  }
  return result.result;
}
