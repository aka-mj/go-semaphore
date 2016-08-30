# go-semaphore

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dangerousHobo/go-semaphore/blob/master/COPYING)

This library provides a wrapper around the C interface for named userspace semaphores on Linux.

## sem_overview(7)

    POSIX semaphores allow processes and threads to synchronize their actions.

       A  semaphore is an integer whose value is never allowed to fall below zero.  Two
       operations can be performed on semaphores: increment the semaphore value by  one
       (sem_post(3));  and  decrement the semaphore value by one (sem_wait(3)).  If the
       value of a semaphore is currently zero, then a sem_wait(3) operation will  block
       until the value becomes greater than zero.

       POSIX semaphores come in two forms: named semaphores and unnamed semaphores.

     Named semaphores
        A named semaphore is identified by a name of the form /somename; that is,
        a null-terminated string of up to NAME_MAX-4 (i.e., 251) characters  con‐
        sisting  of an initial slash, followed by one or more characters, none of
        which are slashes.  Two processes can operate on the same named semaphore
        by passing the same name to sem_open(3).

        The sem_open(3) function creates a new named semaphore or opens an exist‐
        ing named semaphore.  After the semaphore has  been  opened,  it  can  be
        operated  on  using sem_post(3) and sem_wait(3).  When a process has fin‐
        ished using the semaphore, it can use sem_close(3)  to  close  the  sema‐
        phore.   When  all processes have finished using the semaphore, it can be
        removed from the system using sem_unlink(3).

