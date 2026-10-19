const marker = '<!-- tower-defense-preview -->';

export async function updatePreviewComment(github, pr, url, sha, runURL) {
  if (!/^[1-9][0-9]*$/.test(String(pr)) || url !== `https://pr-${pr}.tower-def.pages.dev/` ||
      !/^[a-f0-9]{40}$/.test(sha) || !/^https:\/\/github\.com\/Warsinger\/tower-defense-go\/actions\/runs\/[0-9]+$/.test(runURL)) {
    throw new Error('Invalid preview comment target');
  }
  const body = `${marker}\n## Browser preview\n\n[Play Tower Defense](${url})\n\nRevision: \`${sha}\` · [Build and deployment](${runURL})\n`;
  let existing;
  for (let page = 1; ; page++) {
    const comments = await github(`issues/${pr}/comments?per_page=100&page=${page}`);
    existing = comments.find(comment => comment.user?.login === 'github-actions[bot]' && comment.body?.includes(marker));
    if (existing || comments.length < 100) break;
  }
  if (existing) {
    if (existing.body === body) return existing;
    return github(`issues/comments/${existing.id}`, { method: 'PATCH', body: { body } });
  }
  return github(`issues/${pr}/comments`, { method: 'POST', body: { body } });
}
