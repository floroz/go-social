import {
  useQuery,
  useMutation,
  useQueryClient,
  UseMutationOptions,
} from "@tanstack/react-query";
import PostService from "@/services/postService";
import type {
  ListPostsSuccessResponse,
  CreatePostRequest,
  CreatePostSuccessResponse,
  Post, // Individual post type
  // ApiErrorResponse, // Not used in current implementation, can be added for typed errors
} from "@/types/api";

const POSTS_QUERY_KEY = "posts";

/**
 * Custom hook for fetching posts.
 */
export function usePosts() {
  const { data, isLoading, error, refetch } = useQuery<
    ListPostsSuccessResponse, // This is { data: Post[] }
    Error // Error type for fetching
  >({
    queryKey: [POSTS_QUERY_KEY],
    queryFn: PostService.listPosts,
    // staleTime: 5 * 60 * 1000, // Optional: 5 minutes
  });

  // Extract the actual array of posts from the wrapped response
  const posts: Post[] | undefined = data?.data;

  return {
    posts, // The actual array of Post objects
    isLoading,
    error,
    refetchPosts: refetch,
    // rawResponse: data, // Optionally expose the raw response if needed elsewhere
  };
}

/**
 * Custom hook for creating a new post.
 */
export function useCreatePost(
  options?: UseMutationOptions<
    CreatePostSuccessResponse, // This is { data: Post }
    Error, // Error type for mutation
    CreatePostRequest // Variables type (data sent to mutationFn)
  >
) {
  const queryClient = useQueryClient();

  const { mutate, isPending, error, data } = useMutation<
    CreatePostSuccessResponse,
    Error,
    CreatePostRequest
  >({
    mutationFn: PostService.createPost,
    onSuccess: (newPostResponse, ...rest) => {
      // Invalidate and refetch the posts list to include the new post
      queryClient.invalidateQueries({ queryKey: [POSTS_QUERY_KEY] });
      // Or, for optimistic updates:
      // queryClient.setQueryData([POSTS_QUERY_KEY], (oldData: ListPostsSuccessResponse | undefined) => {
      //   const newPost = newPostResponse.data; // Extract the Post object
      //   return oldData ? { ...oldData, data: [newPost, ...oldData.data] } : { data: [newPost] };
      // });

      // Call original onSuccess if provided
      options?.onSuccess?.(newPostResponse, ...rest);
    },
    onError: options?.onError, // Pass through onError
  });

  return {
    createPost: mutate,
    isCreating: isPending,
    createError: error,
    createdPostData: data, // This is { data: Post }
  };
}

// TODO: Add hooks for updatePost, deletePost, getPostById as needed
