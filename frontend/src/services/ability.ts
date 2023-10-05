import { defineAbility } from "@casl/ability";

export default defineAbility((can, cannot) => {
  can("read", "Post");
  can("update", "Post");
  can("read", "Comment");
  can("update", "Comment");

  cannot('delete', 'Post', { published: true })
});
