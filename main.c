#include <stdio.h>

int main(int const argc, char const *const *argv)
{
   printf("hacked with %d arguments\n", argc);
   for (int i = 0; i < argc; i++)
   {
      printf("%02d: '%s'\n", i, argv[i]);
   }
   return 0;
}
