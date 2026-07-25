#include <stdio.h>

#include "hacked/core/Core.h"

int main(int const argc, char const *argv[])
{
   printf("hacked with %d arguments; sampleReturn: %d\n", argc, sampleReturn(2));
   for (int i = 0; i < argc; i++)
   {
      printf("%02d: '%s'\n", i, argv[i]);
   }
   return 0;
}
