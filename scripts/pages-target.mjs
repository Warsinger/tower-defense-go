export async function deploymentBranch(env, github) {
  const repo = 'Warsinger/tower-defense-go';
  const sha = env.GITHUB_SHA;
  if (env.GITHUB_REPOSITORY !== repo || !/^[a-f0-9]{40}$/.test(sha || '')) throw new Error('Invalid repository/revision');
  if (!['push', 'workflow_dispatch'].includes(env.GITHUB_EVENT_NAME)) throw new Error('PR events cannot publish');
  if (env.GITHUB_REF === 'refs/heads/main') {
    const current = await github('git/ref/heads/main');
    if (current.object.sha !== sha) throw new Error('Refusing to publish a superseded main revision');
    return 'main';
  }
  if (env.GITHUB_EVENT_NAME !== 'workflow_dispatch') throw new Error('Previews require manual dispatch');
  const pr = env.PREVIEW_PR;
  if (!/^[1-9][0-9]*$/.test(pr || '')) throw new Error('Manual preview requires an open PR number');
  const current = await github(`pulls/${pr}`);
  if (current.state !== 'open' || current.base.ref !== 'main' || current.head.repo?.full_name !== repo ||
      current.head.sha !== sha || `refs/heads/${current.head.ref}` !== env.GITHUB_REF) {
    throw new Error('Preview must match the current head of an open same-repository PR targeting main');
  }
  return `pr-${pr}`;
}
